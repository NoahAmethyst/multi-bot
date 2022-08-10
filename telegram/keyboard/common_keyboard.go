package keyboard

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var CommonKeyboard = tgbotapi.ReplyKeyboardMarkup{}

func init() {

	ckb := []tgbotapi.KeyboardButton{{
		Text: ButtonWallet,
	}}
	CommonKeyboard = tgbotapi.NewReplyKeyboard(ckb)
}
