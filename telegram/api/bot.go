package api

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"multi-bot/config"
	"multi-bot/constant"
	"multi-bot/utils/log"
	"net/url"
	"os"
)

type Bot struct {
	*tgbotapi.BotAPI
	BotName string
	Running bool
	Handler func(update tgbotapi.Update, token string)
	BotCmd  map[string]int
	WhInfo  *tgbotapi.WebhookConfig
	TgCmd   []tgbotapi.BotCommand
}

func InitBot(token string) (*Bot, error) {

	if testToken := os.Getenv("TEST_BOT_TOKEN"); testToken != "" {
		token = testToken
	}
	var bot Bot
	b, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Error().Fields(map[string]interface{}{
			"action": "init telegram bot",
			//"token":  token,
			"error": err,
		}).Send()
		panic(map[string]interface{}{
			"action": "init telegram bot",
			//"token":  token,
			"error": err,
		})
		return &bot, err
	}
	bot.BotAPI = b
	botInfo, err := bot.GetMe()
	if err != nil {
		log.Error().Fields(map[string]interface{}{
			"action": "get bot name",
			//"token":  token,
			"error": err,
		}).Send()
		panic(err)
	}
	bot.BotName = botInfo.UserName

	return &bot, nil
}

func (b *Bot) Start() {
	if config.EnvIsDev() {
		b.StartDev()
	} else {
		b.StartPrd()
	}
}

func (b *Bot) StartPrd() {

	var updates tgbotapi.UpdatesChannel
	if config.UseWebHook() {
		whUrl, err := url.Parse(os.Getenv(constant.TG_WH_URL) + b.BotName + "/")
		if err != nil {
			log.Error().Msgf("invalid webhook addr %s, error %s", os.Getenv(constant.TG_WH_URL), err.Error())
			return
		}
		tgListen := whUrl.RequestURI()

		log.Info().Msgf("tg listen uri %s", tgListen)

		updates = b.ListenForWebhook(tgListen)

		whAddr := os.Getenv(constant.TG_WH_URL) + b.BotName + "/"
		if whAddr == "" {
			log.Error().Fields(map[string]interface{}{
				"action":   "start telegram bot",
				"error":    fmt.Errorf("empty tg webhook url"),
				"bot name": b.BotName,
			}).Send()
			return
		}

		//time.Sleep(time.Second * 10)
		wh, err := tgbotapi.NewWebhook(whAddr)
		if err != nil {
			log.Error().Fields(map[string]interface{}{
				"action":   "start telegram bot",
				"error":    err,
				"bot name": b.BotName,
			}).Send()
			return
		}

		log.Info().Msgf("tg bot webhook start at %s", whAddr)

		resp, err := b.Request(wh)
		if err != nil {
			log.Panic().Fields(map[string]interface{}{
				"action":   "start telegram bot",
				"error":    err,
				"bot name": b.BotName,
			}).Send()
			return
		} else {
			log.Info().Fields(map[string]interface{}{
				"action":   "start telegram bot",
				"resp":     resp,
				"bot name": b.BotName,
			}).Send()
		}
	} else {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		//b.Debug = true
		updates = b.GetUpdatesChan(u)
		dw := tgbotapi.DeleteWebhookConfig{DropPendingUpdates: true}
		_, err := b.Request(dw)
		if err != nil {
			log.Error().Msgf("delete webhook error %s", err.Error())
			return
		}
	}
	log.Info().Fields(map[string]interface{}{
		"action": "start telegram bot",
		"name":   b.BotAPI.Self.UserName,
	}).Send()

	for update := range updates {
		update := update
		go b.Handler(update, b.BotName)
	}
}

func (b *Bot) StartDev() {
	var updates tgbotapi.UpdatesChannel
	if config.UseWebHook() {
		var listenUrl string
		listenUrl = "/tg-webhook/"
		whUrl, err := url.Parse(listenUrl + b.BotName + "/")
		if err != nil {
			log.Error().Msgf("invalid webhook addr %s, error %s", listenUrl, err.Error())
			return
		}
		tgListen := whUrl.RequestURI()

		log.Info().Msgf("tg listen uri %s", tgListen)

		updates = b.ListenForWebhook(tgListen)

		whAddr := os.Getenv(constant.TG_WH_URL) + "/" + b.BotName + "/"

		if whAddr == "" {
			log.Error().Fields(map[string]interface{}{
				"action":   "start telegram bot",
				"error":    fmt.Errorf("empty tg webhook url"),
				"bot name": b.BotName,
			}).Send()
			return
		}

		//time.Sleep(time.Second * 10)
		wh, err := tgbotapi.NewWebhook(whAddr)
		if err != nil {
			webhook, _ := b.GetWebhookInfo()
			log.Error().Fields(map[string]interface{}{
				"action":   "start telegram bot",
				"error":    err,
				"bot name": b.BotName,
				"webhook":  webhook,
			}).Send()
			return
		}

		log.Info().Msgf("tg bot webhook start at %s", whAddr)

		resp, err := b.Request(wh)
		if err != nil {
			log.Panic().Fields(map[string]interface{}{
				"action":   "start telegram bot",
				"error":    err,
				"bot name": b.BotName,
			}).Send()
			return
		} else {
			log.Info().Fields(map[string]interface{}{
				"action":   "start telegram bot",
				"resp":     resp,
				"bot name": b.BotName,
			}).Send()
		}
	} else {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		//b.Debug = true
		updates = b.GetUpdatesChan(u)
		dw := tgbotapi.DeleteWebhookConfig{DropPendingUpdates: true}
		_, err := b.Request(dw)
		if err != nil {
			log.Error().Msgf("delete webhook error %s", err.Error())
			return
		}
	}

	log.Info().Fields(map[string]interface{}{
		"action": "start telegram bot",
		"name":   b.BotAPI.Self.UserName,
	}).Send()

	for update := range updates {
		update := update
		go b.Handler(update, b.BotName)
	}
}
