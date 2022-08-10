package manager

import (
	"github.com/bwmarrin/discordgo"
	"multi-bot/discord/cmd"
	"multi-bot/discord/message_component_handler"
	"multi-bot/utils/log"
)

func NewDiscordCmdHandler(app *Application) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		HandleCmd(app, s, i)
	}
}

func HandleCmd(app *Application, s *discordgo.Session, i *discordgo.InteractionCreate) {

	handleErr := handleCmd(app, s, i)
	if handleErr != nil {
		log.Error().Fields(map[string]interface{}{"action": "handle cmd", "error": handleErr.Error()}).Send()
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: handleErr.Error(),
			},
		})
		if err != nil {
			log.Error().Fields(map[string]interface{}{"action": "reply error", "error": err.Error()}).Send()
		}
		//var userId string
		//if i.Member != nil {
		//	userId = i.Member.User.ID
		//} else {
		//	userId = i.User.ID
		//}
		//channel, err := s.UserChannelCreate(userId)
		//if err != nil {
		//	log.Error().Fields(map[string]interface{}{"action": "user channel create error", "error": err.Error(), "userId": userId, "cmdError": handleErr.Error()}).Send()
		//	return
		//}
		//_, err = s.ChannelMessageSend(channel.ID, handleErr.Error())
		//if err != nil {
		//	log.Error().Fields(map[string]interface{}{"action": "user channel msg error", "error": err.Error(), "userId": userId, "cmdError": handleErr.Error()}).Send()
		//}
	} else {
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Operation Success",
			},
		}); err != nil {
			log.Error().Msgf("discord interaction respond error:%s", err.Error())
		}

	}
}

func handleCmd(app *Application, s *discordgo.Session, i *discordgo.InteractionCreate) error {

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		handle := cmd.GetCmdHandler(i.ApplicationCommandData().Name)
		if handle != nil {
			return handle.Handler(*app.AppContext, s, i)
		}
	case discordgo.InteractionMessageComponent:
		h := message_component_handler.GetMsgComponentHandler(i.Interaction.MessageComponentData().CustomID)
		if h != nil {
			return h(*app.AppContext, s, i)
		}
	}
	log.Info().Fields(map[string]interface{}{"action": "unknown interaction", "i": i.Interaction}).Send()
	return nil
}
