package message_handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"multi-bot/app_context"
)

//Todo
//Filter spam message
func FilterMessage(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, message *tgbotapi.Message) {

}
