package message_handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"multi-bot/app_context"
	"multi-bot/telegram/keyboard"
	"multi-bot/utils/log"
)

var buttonCmdMap map[string]bool

func CheckIsButton(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, message *tgbotapi.Message) {
	if v, ok := buttonCmdMap[message.Text]; !ok || !v {
		return
	}

	if err := buttonProxy(appCtx, fromChat, sender, message.Text); err != nil {
		_, _ = appCtx.TgBotApi.SendMsg(fromChat, err.Error(), nil, false, false)
	}

}

func buttonProxy(appCtx *app_context.AppContext, fromChat *tgbotapi.Chat, sender *tgbotapi.User, button string) error {
	switch button {

	default:
		log.Warn().Msgf("unhandled button :%s", button)
		return nil
	}
}

func init() {
	buttonCmdMap = map[string]bool{
		keyboard.ButtonWallet: true,
	}
}
