package cmd

import (
	"github.com/bwmarrin/discordgo"
	"multi-bot/config"
	"multi-bot/discord/cmdhandler"
)

var handleMap = map[string]*cmdhandler.DiscordCmdHandler{
	"create_quiz": &cmdhandler.DiscordCmdHandler{Handler: cmdhandler.CreateQuizCallback, ApplicationCommand: &discordgo.ApplicationCommand{
		ID:                       "",
		ApplicationID:            "",
		Version:                  "",
		Type:                     0,
		Name:                     "create_quiz",
		NameLocalizations:        nil,
		DefaultPermission:        nil,
		Description:              "quiz dev test cmd",
		DescriptionLocalizations: nil,
	}},
}

func init() {
	if config.EnvIsDev() {
		handleMap["create_quiz_test"] = &cmdhandler.DiscordCmdHandler{Handler: cmdhandler.CreateQuiz, ApplicationCommand: &discordgo.ApplicationCommand{
			ID:                       "",
			ApplicationID:            "",
			Version:                  "",
			Type:                     0,
			Name:                     "create_quiz_test",
			NameLocalizations:        nil,
			DefaultPermission:        nil,
			Description:              "create_quiz_test",
			DescriptionLocalizations: nil,
			Options:                  nil,
		}}
		handleMap["stop_quiz_test"] = &cmdhandler.DiscordCmdHandler{Handler: cmdhandler.StopQuiz, ApplicationCommand: &discordgo.ApplicationCommand{
			ID:                       "",
			ApplicationID:            "",
			Version:                  "",
			Type:                     0,
			Name:                     "stop_quiz_test",
			NameLocalizations:        nil,
			DefaultPermission:        nil,
			Description:              "stop_quiz_test",
			DescriptionLocalizations: nil,
			Options:                  nil,
		},
		}
	}
}

func GetDiscordCmdList() []*discordgo.ApplicationCommand {
	resp := make([]*discordgo.ApplicationCommand, 0)
	for _, v := range handleMap {
		resp = append(resp, v.ApplicationCommand)
	}
	return resp
}
func GetCmdHandler(cmdId string) *cmdhandler.DiscordCmdHandler {
	return handleMap[cmdId]
}
