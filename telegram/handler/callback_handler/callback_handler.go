package callback_handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"multi-bot/app_context"
	"multi-bot/utils/log"
)

type CallbackHandler struct {
	Handlers []func(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, callback *tgbotapi.CallbackQuery)
}

func NewCallbackHandler(handlers ...func(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, callback *tgbotapi.CallbackQuery)) *CallbackHandler {
	return &CallbackHandler{
		Handlers: handlers,
	}
}

func (h *CallbackHandler) AddHandler(handler func(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, callback *tgbotapi.CallbackQuery)) *CallbackHandler {
	h.Handlers = append(h.Handlers, handler)
	return h
}

func (h *CallbackHandler) Handle(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, callback *tgbotapi.CallbackQuery) {
	if callback == nil {
		return
	}
	for _, handler := range h.Handlers {
		handler(appCtx, fromChat, sender, callback)
	}
}

func Logger(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, callback *tgbotapi.CallbackQuery) {
	if callback == nil {
		return
	}
	log.Info().Fields(map[string]interface{}{
		"app":      appCtx,
		"action":   "receive callback",
		"chat":     fromChat,
		"sender":   sender,
		"callback": callback,
	}).Send()
}
