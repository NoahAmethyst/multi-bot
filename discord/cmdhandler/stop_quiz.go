package cmdhandler

import (
	"multi-bot/app_context"

	"github.com/bwmarrin/discordgo"
)

func StopQuiz(appCtx app_context.AppContext, s *discordgo.Session, i *discordgo.InteractionCreate) error {

	return nil
	//return quizClient.StopQuiz(i.ChannelID)
}
