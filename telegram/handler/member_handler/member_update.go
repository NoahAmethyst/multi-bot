package member_handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"multi-bot/app_context"
)

const (
	statusLeft = "left"
)

func MemberUpdate(appCtx *app_context.AppContext, _ *tgbotapi.Chat, _ *tgbotapi.User, memberUpdate *tgbotapi.ChatMemberUpdated) {

	BotRemove(appCtx, memberUpdate)

}

func BotRemove(appCtx *app_context.AppContext, memberUpdate *tgbotapi.ChatMemberUpdated) {
	if !memberUpdate.OldChatMember.User.IsBot {
		return
	}
	if memberUpdate.NewChatMember.Status != statusLeft {
		return
	}
	if memberUpdate.OldChatMember.User.UserName == appCtx.TgBotApi.BotName {
		appCtx.DelGroup(memberUpdate.Chat.ID)
	}
}
