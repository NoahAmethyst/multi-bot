package poll_handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"multi-bot/app_context"
	"multi-bot/utils/log"
)

type PollAnswerHandler struct {
	Handlers []func(appCtx *app_context.AppContext, pollAnswer *tgbotapi.PollAnswer)
}

func NewPollAnswerHandler(handlers ...func(appCtx *app_context.AppContext, pollAnswer *tgbotapi.PollAnswer)) *PollAnswerHandler {
	return &PollAnswerHandler{
		Handlers: handlers,
	}
}

func (h *PollAnswerHandler) AddHandler(handler func(appCtx *app_context.AppContext, pollAnswer *tgbotapi.PollAnswer)) *PollAnswerHandler {
	h.Handlers = append(h.Handlers, handler)
	return h
}

func (h *PollAnswerHandler) Handle(appCtx *app_context.AppContext, pollAnswer *tgbotapi.PollAnswer) {
	if pollAnswer == nil {
		return
	}
	for _, handler := range h.Handlers {
		handler(appCtx, pollAnswer)
	}
}

func LoggerPollAnswer(appCtx *app_context.AppContext, pollAnswer *tgbotapi.PollAnswer) {
	if pollAnswer == nil {
		return
	}
	log.Info().Fields(map[string]interface{}{
		"app":        appCtx,
		"action":     "receive poll answer",
		"pollAnswer": pollAnswer,
	}).Send()
}
