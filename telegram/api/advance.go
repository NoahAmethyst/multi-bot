package api

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func (b *Bot) DelDeadMsg(msg *tgbotapi.Message, deadline time.Duration) {
	if msg == nil || msg.Chat == nil || deadline == 0 {
		return
	}
	go func() {
		time.Sleep(deadline)
		_ = b.DelMsg(msg)
	}()
}

func (b *Bot) PinTemporaryMsg(chatId int64, messageId int, deadline time.Duration) {
	if err := b.PinMsg(chatId, messageId); err != nil || deadline == 0 {
		return
	}
	go func() {
		time.Sleep(deadline)
		_ = b.UnPinMsg(chatId, messageId)
	}()
}

func (b *Bot) BatchSendMsg() {
	//Todo @mazhonghao
}
