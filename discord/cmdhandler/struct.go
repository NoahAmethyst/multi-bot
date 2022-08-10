package cmdhandler

import (
	"github.com/bwmarrin/discordgo"
	"multi-bot/app_context"
)

type DiscordCmdHandler struct {
	ApplicationCommand *discordgo.ApplicationCommand                                                                   `json:"detail"`
	Handler            func(appCtx app_context.AppContext, s *discordgo.Session, i *discordgo.InteractionCreate) error `json:"handler"`
}
