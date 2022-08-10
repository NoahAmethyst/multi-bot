package message_handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"multi-bot/app_context"
	"multi-bot/utils/log"
)

type MessageHandler struct {
	Handlers []func(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, message *tgbotapi.Message)
}

func NewMessageHandler(handlers ...func(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, message *tgbotapi.Message)) *MessageHandler {
	return &MessageHandler{
		Handlers: handlers,
	}
}

func (h *MessageHandler) AddHandler(handler func(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, message *tgbotapi.Message)) *MessageHandler {
	h.Handlers = append(h.Handlers, handler)
	return h
}

func (h *MessageHandler) Handle(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, message *tgbotapi.Message) {
	if message == nil {
		return
	}
	for _, handler := range h.Handlers {
		handler(appCtx, fromChat, sender, message)
	}
}

func Logger(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, message *tgbotapi.Message) {
	if message == nil {
		return
	}
	log.Info().Fields(map[string]interface{}{
		"app":     appCtx,
		"action":  "receive message",
		"chat":    fromChat,
		"sender":  sender,
		"message": message,
	}).Send()
}
