package manager

import (
	"github.com/bwmarrin/discordgo"
	"google.golang.org/grpc"
	"multi-bot/app_context"
	"multi-bot/cluster/grpc_client"
	"multi-bot/constant"
	discord_api "multi-bot/discord/api"
	"multi-bot/discord/cmd"
	"multi-bot/entity/entity_pb/app_pb"
	tgapi "multi-bot/telegram/api"
	"multi-bot/telegram/handler/member_handler"
	"multi-bot/telegram/handler/message_handler"
	"multi-bot/telegram/handler/poll_handler"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"multi-bot/db/applicationdb"
	"multi-bot/utils/log"
)

var applications map[string]Application

type Application struct {
	*app_context.AppContext
}

type handleRecord struct {
	sync.RWMutex
	handleRecords map[int]bool
}

var handleRecords handleRecord

func Handler(app Application, update *tgbotapi.Update, botName string) {
	//log.Info().Fields(map[string]interface{}{
	//	"update": update,
	//}).Send()

	if handleRecords.CheckIsHandle(update.UpdateID) {
		return
	}

	if botName != app.TgBotApi.BotName {
		log.Warn().Msgf("register bot:%s income bot:%s", app.TgBotApi.BotName, botName)
		return
	}

	handleRecords.RecordHandle(update.UpdateID)

	message_handler.NewMessageHandler(message_handler.FilterMessage).
		AddHandler(message_handler.MemberJoin).
		AddHandler(message_handler.MemberLeft).
		AddHandler(message_handler.NewChatTitle).
		AddHandler(message_handler.CheckIsButton).
		AddHandler(message_handler.CmdHandler).
		Handle(app.AppContext, update.FromChat(), update.SentFrom(), update.Message)

	poll_handler.NewPollAnswerHandler(poll_handler.Analyse).
		Handle(app.AppContext, update.PollAnswer)

	member_handler.NewMemberUpdateHandler(member_handler.MemberUpdate).Handle(app.AppContext, update.FromChat(), update.SentFrom(), update.MyChatMember)
}

func GetApplication(appId string) (Application, bool) {
	v, ok := applications[appId]
	return v, ok
}

func InitAllApps() {
	apps, err := applicationdb.GetAllApps()

	if err != nil {
		log.Error().Msgf("get all apps error: %s", err)
		return
	}

	for _, app := range apps {
		var connection *grpc.ClientConn
		var uri string
		if app.SvcType == app_pb.SvcType_Grpc {
			connection, err = grpc_client.StartConnection(app.SvcUrl)
			if err != nil {
				log.Error().Msgf("start connection error %s", err)
				continue
			}
		} else {
			uri = app.SvcUrl
		}

		//set default group
		groups, _ := applicationdb.GetTgApplicationGroups(app.AppId, constant.PlatTG)

		var application Application
		if v, ok := applications[app.AppId]; !ok {
			application = Application{
				AppContext: &app_context.AppContext{
					Id:     app.AppId,
					Name:   app.Name,
					Conn:   connection,
					Url:    uri,
					Groups: groups,
				},
			}
		} else {
			application = v
		}

		switch app.Type {
		case constant.PlatTG:
			//set bot
			tBot, err := tgapi.InitBot(app.BotToken)
			if err != nil {
				log.Error().Msgf("app %s init bot error:%s", app.Name, err)
				continue
			}
			var cmdList []tgbotapi.BotCommand
			for cmdStr, cmdDesc := range app.CmdConfig.CmdConfig {
				cmdList = append(cmdList, tgbotapi.BotCommand{
					Command:     cmdStr,
					Description: cmdDesc,
				})
			}
			tBot.RegisterCmdList(cmdList)
			application.SetApplicationBot(tBot)
		case constant.PlatDiscord:
			dbot, err := discord_api.InitBot(app.BotToken, app.AppId)
			if err != nil {
				log.Error().Msgf("app %s init bot error:%s", app.Name, err)
				continue
			}
			application.DiscordBotApi = dbot
		}

		applications[app.AppId] = application

	}
}

func (r *handleRecord) RecordHandle(updateId int) {
	r.Lock()
	defer r.Unlock()
	r.handleRecords[updateId] = true
}

func (r *handleRecord) CheckIsHandle(updateId int) bool {
	r.RLock()
	defer r.RUnlock()
	return r.handleRecords[updateId]
}

func init() {
	applications = map[string]Application{}
	handleRecords = handleRecord{
		RWMutex:       sync.RWMutex{},
		handleRecords: map[int]bool{},
	}
}
func StartAllBots() {
	InitAllApps()
	for _, app := range applications {
		app := app
		app.startBot()
	}
}

func (a *Application) startBot() {

	var hasBot bool

	if a.TgBotApi != nil {
		hasBot = true
		a.TgBotApi.SetHandler(teleBotAppHandler(*a))
		log.Info().Msgf("app %s start bot", a.Name)
		go a.TgBotApi.Start()
	}
	if a.DiscordBotApi != nil {
		hasBot = true
		go func(app *Application) {
			log.Info().Fields(map[string]interface{}{"action": "starting init discord bot"}).Send()
			a.DiscordBotApi.PutHandler(NewDiscordCmdHandler(a))
			a.DiscordBotApi.PutHandler(func(s *discordgo.Session, i *discordgo.GuildCreate) {
				log.Info().Fields(map[string]interface{}{"action": "join server", "appId": app.Id, "serverName": i.Name}).Send()
				//for _, ch := range i.Channels {
				//	if ch.Type == discordgo.ChannelTypeGuildText {
				//		chId, err := strconv.ParseInt(ch.ID, 10, 64)
				//		if err != nil {
				//			log.Error().Fields(map[string]interface{}{"action": "invalid channel id", "error": err.Error(), "id": chId}).Send()
				//		}
				//	}
				//}
			})
			a.DiscordBotApi.SetCmds(cmd.GetDiscordCmdList())
			err := app.DiscordBotApi.Start()
			if err != nil {
				log.Error().Fields(map[string]interface{}{"action": "start discord bot error", "error": err.Error()}).Send()
			} else {
				log.Info().Msgf("app %s start bot", app.Name)
			}
		}(a)
	}

	if !hasBot {
		log.Warn().Msgf(constant.AppBotNotSetDetail, a.Name)
		return
	}

	//handler := a.teleBotAppHandler()

}
