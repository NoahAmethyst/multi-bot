package manager

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"multi-bot/telegram/handler/callback_handler"
	"multi-bot/telegram/handler/message_handler"
	"multi-bot/telegram/handler/poll_handler"
)

func teleBotAppHandler(app Application) func(update tgbotapi.Update, botName string) {
	return func(update tgbotapi.Update, botName string) {
		Handler(app, &update, botName)
	}
}

func defaultHandler() func(update tgbotapi.Update) {
	return func(update tgbotapi.Update) {
		message_handler.NewMessageHandler(message_handler.Logger).
			Handle(nil, update.FromChat(), update.SentFrom(), update.Message)

		poll_handler.NewPollAnswerHandler(poll_handler.LoggerPollAnswer).
			Handle(nil, update.PollAnswer)

		callback_handler.NewCallbackHandler(callback_handler.Logger).
			Handle(nil, update.FromChat(), update.SentFrom(), update.CallbackQuery)
	}
}
