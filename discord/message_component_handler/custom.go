package message_component_handler

import (
	"github.com/bwmarrin/discordgo"
	"multi-bot/app_context"
)

var componentMap = map[string]func(appCtx app_context.AppContext, s *discordgo.Session, i *discordgo.InteractionCreate) error{}

func init() {
	componentMap = map[string]func(appCtx app_context.AppContext, s *discordgo.Session, i *discordgo.InteractionCreate) error{}
}

func GetMsgComponentHandler(customId string) func(appCtx app_context.AppContext, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return componentMap[customId]
}
