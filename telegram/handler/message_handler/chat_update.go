package message_handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"multi-bot/app_context"
	"time"
)

//Todo
//Filter spam message
func MemberLeft(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, message *tgbotapi.Message) {
	if message.LeftChatMember == nil {
		return
	}
	appCtx.TgBotApi.DelDeadMsg(message, 5*time.Second)
}

func MemberJoin(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, message *tgbotapi.Message) {
	if message.NewChatMembers == nil || len(message.NewChatMembers) == 0 {
		return
	}

	BotJoin(appCtx, fromChat, message.NewChatMembers[0])

	appCtx.TgBotApi.DelDeadMsg(message, 5*time.Second)
}

func NewChatTitle(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, message *tgbotapi.Message) {
	if len(message.NewChatTitle) == 0 {
		return
	}

	if appCtx.CheckInGroup(fromChat.ID) {
		appCtx.ChangeGroup(fromChat.ID, fromChat.Title)
	}

}

func BotJoin(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, newChatMember tgbotapi.User) {
	if !newChatMember.IsBot {
		return
	}
	//Todo not app bot
	if newChatMember.UserName != appCtx.TgBotApi.BotName {

	} else {
		if !appCtx.CheckInGroup(fromChat.ID) {
			appCtx.AddGroup(fromChat.ID, fromChat.Title)
		}
	}
}
