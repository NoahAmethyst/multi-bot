package api

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"multi-bot/constant"
	"multi-bot/utils/dingding"

	"multi-bot/utils/log"
	"os"
	"time"
)

type Bot struct {
	appId string
	name  string
	*discordgo.Session
	//commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	//commands        map[string]discordgo.ApplicationCommand
	commands []*discordgo.ApplicationCommand
}

func (b *Bot) GetName() string {
	return b.name
}

func InitBot(token string, appId string) (*Bot, error) {

	if testToken := os.Getenv("DISCORD_BOT_TOKEN"); testToken != "" {
		token = testToken
	}

	var bot Bot
	c, err := discordgo.New("Bot " + token)

	if err != nil {
		log.Error().Fields(map[string]interface{}{
			"action": "init discord bot",
			"token":  token,
			"error":  err,
		}).Send()
		return &bot, err
	}

	//c.Client.Timeout = time.Second * 10
	bot.Session = c
	bot.appId = appId

	log.Info().Fields(map[string]interface{}{"action": "init discord bot success"}).Send()

	return &bot, nil
}

//
//func (b *Bot) SetCmdHandler(cmdS string, desc string, options []*discordgo.ApplicationCommandOption, f func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
//
//	if b.commandHandlers == nil {
//		b.commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
//	}
//	if b.commands == nil {
//		b.commands = map[string]discordgo.ApplicationCommand{}
//	}
//	b.commands[cmdS] = discordgo.ApplicationCommand{
//		Name:        cmdS,
//		Description: desc,
//		Options:     options,
//	}
//
//	b.commandHandlers[cmdS] = f
//
//}

func (b *Bot) SetCmds(cmds []*discordgo.ApplicationCommand) {
	b.commands = cmds
}

func (b *Bot) RegisterCmds() error {
	existCmdList, err := b.ApplicationCommands(b.State.User.ID, "")
	if err != nil {
		log.Error().Fields(map[string]interface{}{"action": "get exist cmd error", "error": err.Error()}).Send()
	}

	if len(b.commands) > len(existCmdList) || os.Getenv("RESET_DISCORD_CMD") == "1" {
		for _, v := range b.commands {
			_, err = b.ApplicationCommandCreate(b.State.User.ID, "", v)
			log.Info().Fields(map[string]interface{}{"action": "register cmd", "cmdId": v.Name}).Send()
			if err != nil {
				log.Error().Fields(map[string]interface{}{"action": "create cmd error", "error": err.Error(), "cmd": v.Name}).Send()
				return err
			}
		}

	}
	return nil
}

func (b *Bot) PutHandler(f interface{}) {
	b.AddHandler(f)
}

func (b *Bot) Start() error {
	log.Info().Fields(map[string]interface{}{"action": "start init discord bot"}).Send()
	// In this example, we only care about receiving message events.
	b.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	//b.AddHandler(manager.HandleCmd)
	b.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
	})

	// Open a websocket connection to Discord and begin listening.
	err := b.Open()

	if err != nil {
		retryMax := 3
		for i := 1; i <= retryMax; i++ {
			time.Sleep(time.Duration(i*15) * time.Second)
			if err = b.Open(); err == nil {
				break
			} else {
				if i == retryMax {
					content := fmt.Sprintf(`### Discord 机器人告警,链接重试 %d 次失败，请查看`, i)
					robot := dingding.NewRobot(constant.DINGDING_TOKEN, constant.DINGDING_SECRET, constant.DINGDING_APP_KEY, constant.DINGDING_SECRET_KEY)
					if err := robot.SendMarkdownMessage("##Survival Bot Daily Report", content, []string{"+86-18100170551", "+86-15251724436"}, true); err != nil {
						log.Error().Msgf("send dingding msg error:%s", err)
					}
					log.Error().Fields(map[string]interface{}{
						"action": "start discord bot",
						"error":  err,
					}).Send()
					return err
				}
			}
		}
	}
	b.name = b.Session.State.User.Username
	if err := b.RegisterCmds(); err != nil {
		log.Error().Fields(map[string]interface{}{"action": "register cmd error", "error": err.Error()}).Send()
		return err
	}

	log.Info().Fields(map[string]interface{}{"action": "open discord connection success"}).Send()

	log.Info().Fields(map[string]interface{}{
		"action": "start discord bot",
		"name":   b.State.User.Username,
	}).Send()

	return nil
}
