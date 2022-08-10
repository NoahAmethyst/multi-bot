package message_handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"multi-bot/app_context"
	"multi-bot/utils/log"
	"strings"
)

//Todo
//Filter spam message
func CmdHandler(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, message *tgbotapi.Message) {
	if !message.IsCommand() {
		return
	}
	//check cmd is send to call this bot or not
	if len(message.Text) > 0 && message.Text[0:1] == "/" {
		//if cmd is not call this bot then ignore it
		if !message.Chat.IsPrivate() && strings.Contains(message.Text, "@") {
			if !strings.Contains(message.Text, appCtx.TgBotApi.BotName) {
				return
			}
		}
	}

	if err := cmdProxy(appCtx, fromChat, sender, message.Command()); err != nil {
		_, _ = appCtx.TgBotApi.SendMsg(fromChat, err.Error(), nil, false, false)
	}
}

//Notice
//New command should router here
func cmdProxy(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, cmdStr string) error {
	switch cmdStr {

	default:
		log.Warn().Msgf("no such command:%s", cmdStr)
	}
	return nil
}
