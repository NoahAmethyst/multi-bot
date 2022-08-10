package member_handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"multi-bot/app_context"
	"multi-bot/utils/log"
)

type MemberUpdateHandler struct {
	Handlers []func(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, memberUpdate *tgbotapi.ChatMemberUpdated)
}

func NewMemberUpdateHandler(handlers ...func(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, memberUpdate *tgbotapi.ChatMemberUpdated)) *MemberUpdateHandler {
	return &MemberUpdateHandler{
		Handlers: handlers,
	}
}

func (h *MemberUpdateHandler) AddHandler(handler func(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, memberUpdate *tgbotapi.ChatMemberUpdated)) *MemberUpdateHandler {
	h.Handlers = append(h.Handlers, handler)
	return h
}

func (h *MemberUpdateHandler) Handle(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, memberUpdate *tgbotapi.ChatMemberUpdated) {
	if memberUpdate == nil {
		return
	}
	for _, handler := range h.Handlers {
		handler(appCtx, fromChat, sender, memberUpdate)
	}
}

func Logger(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, memberUpdate *tgbotapi.ChatMemberUpdated) {
	if memberUpdate == nil {
		return
	}
	log.Info().Fields(map[string]interface{}{
		"app":          appCtx,
		"action":       "receive message",
		"chat":         fromChat,
		"sender":       sender,
		"memberUpdate": memberUpdate,
	}).Send()
}
