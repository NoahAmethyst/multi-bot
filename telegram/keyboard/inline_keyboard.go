package keyboard

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"multi-bot/constant"
	"multi-bot/telegram/api"
	"time"
)

//forward keyboard need to be delete after some minutes
func NewForwardPrivateKeyBoard(text string, bot *api.Bot) (*tgbotapi.InlineKeyboardMarkup, time.Duration) {
	km := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonURL(text, fmt.Sprintf("https://t.me/%s", bot.BotName))})
	return &km, constant.COMMON_KEYBOARD_DEADLINE
}
