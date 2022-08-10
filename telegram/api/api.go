package api

import (
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"multi-bot/utils/log"
	"multi-bot/utils/obj_util"
)

const (
	PollQuiz    = "quiz"
	PollRegular = "regular"
)

func (b *Bot) SetHandler(f func(update tgbotapi.Update, botName string)) {
	b.Handler = f
}

func (b *Bot) RegisterCmdList(cmdList []tgbotapi.BotCommand) {

	smcConfig := tgbotapi.SetMyCommandsConfig{Commands: cmdList}
	if _, err := b.Request(smcConfig); err != nil {
		log.Error().Msgf("bot  name(%s) register failed %s", b.BotName, err.Error())
	} else {
		log.Info().Fields(map[string]interface{}{
			"action":  "register cmd",
			"botName": b.BotName,
			"cmd":     smcConfig.Commands,
		}).Send()
	}
}

func (b *Bot) SendMsg(chat *tgbotapi.Chat, content string, ikm interface{}, markdownParse bool, disableWebPreview bool) (tgbotapi.Message, error) {

	msg := tgbotapi.NewMessage(chat.ID, content)
	msg.DisableWebPagePreview = disableWebPreview
	if ikm != nil {
		msg.ReplyMarkup = ikm
	}
	if markdownParse {
		msg.ParseMode = tgbotapi.ModeMarkdownV2
	}
	m, err := b.Send(msg)
	if err != nil {
		log.Error().Fields(map[string]interface{}{
			"action":  "telegram bot send message",
			"token":   b.Token,
			"chat":    chat.Title,
			"content": content,
			"error":   err,
		}).Send()
	}
	return m, err
}

func (b *Bot) SendPhoto(chat *tgbotapi.Chat, photoUrl string, ikm interface{}) (tgbotapi.Message, error) {
	photoConfig := tgbotapi.NewPhoto(chat.ID, tgbotapi.FileURL(photoUrl))
	photoConfig.ReplyMarkup = ikm
	m, err := b.Send(photoConfig)
	if err != nil {
		log.Error().Fields(map[string]interface{}{
			"action": "telegram bot send photo",
			"token":  b.Token,
			"chat":   chat.Title,
			"photo":  photoUrl,
			"error":  err,
		}).Send()
	}
	return m, err
}

func (b *Bot) DelMsg(msg *tgbotapi.Message) error {
	res, err := b.Request(tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID))
	if err != nil {
		log.Error().Fields(map[string]interface{}{
			"action": "telegram bot delete message",
			"name":   b.BotName,
			"chat":   msg.Chat.Title,
			"msg":    msg.Text,
			"error":  err,
		}).Send()
		return err
	}
	if !res.Ok {
		log.Error().Fields(map[string]interface{}{
			"action": "telegram bot delete message",
			"name":   b.BotName,
			"chat":   msg.Chat.Title,
			"msg":    msg.Text,
			"error":  res.Description,
		}).Send()
		err = errors.New(res.Description)
	}
	return err
}

func (b *Bot) ReplyMsg(msg *tgbotapi.Message, content string, ikm interface{}, markdownParse bool, disableWebPreview bool) (tgbotapi.Message, error) {
	reply := tgbotapi.NewMessage(msg.Chat.ID, content)
	reply.DisableWebPagePreview = disableWebPreview

	if ikm != nil {
		reply.ReplyMarkup = ikm
	}

	if markdownParse {
		reply.ParseMode = tgbotapi.ModeMarkdownV2
	}

	reply.ReplyToMessageID = msg.MessageID
	m, err := b.Send(reply)

	if err != nil {
		log.Error().Fields(map[string]interface{}{
			"action":  "telegram bot reply message",
			"token":   b.Token,
			"chat":    msg.Chat.Title,
			"content": content,
			"error":   err,
		}).Send()
	}
	return m, err
}

func (b *Bot) EditMsg(chatId int64, messageId int, content string, ikm *tgbotapi.InlineKeyboardMarkup, markdownParse bool, disableWebPreview bool) error {
	var editMsg tgbotapi.EditMessageTextConfig
	if ikm != nil {
		editMsg = tgbotapi.NewEditMessageTextAndMarkup(chatId, messageId, content, *ikm)
	} else {
		editMsg = tgbotapi.NewEditMessageText(chatId, messageId, content)
	}
	editMsg.DisableWebPagePreview = disableWebPreview
	if markdownParse {
		editMsg.ParseMode = tgbotapi.ModeMarkdownV2
	}
	_, err := b.Request(editMsg)
	if err != nil {
		log.Error().Fields(map[string]interface{}{
			"action": "edit msg",
			"error":  err.Error(),
			"chatId": chatId,
			"msgId":  messageId,
		}).Send()
	}
	return err
}

func (b *Bot) SendQuizPoll(chatId int64, question string, options []string, correctIndex int64, needAnonymous bool) (*tgbotapi.Poll, int64, error) {
	pollCfg := tgbotapi.NewPoll(chatId, question)
	pollCfg.Options = options
	pollCfg.Type = PollQuiz
	pollCfg.Explanation = fmt.Sprintf("correct answer is %d", correctIndex+1)
	pollCfg.CorrectOptionID = correctIndex
	pollCfg.IsAnonymous = needAnonymous
	pollCallback := &tgbotapi.Poll{}
	var msgId int64
	resp, err := b.Request(pollCfg)
	if err != nil {
		log.Error().Fields(map[string]interface{}{
			"action": "send quiz poll",
			"error":  err,
		}).Send()
	} else {
		if !resp.Ok {
			log.Error().Fields(map[string]interface{}{
				"action":       "send quiz poll",
				"request resp": resp,
			}).Send()
			err = errors.New(resp.Description)
		} else {
			var respData map[string]interface{}
			var bytes []byte
			if bytes, err = resp.Result.MarshalJSON(); err == nil {
				if err = json.Unmarshal(bytes, &respData); err != nil {
					log.Error().Fields(map[string]interface{}{
						"action":     "send quiz poll",
						"sub action": "unmarshal request json",
						"error":      err,
					}).Send()
				} else {
					if err = obj_util.MapToStruct(respData["poll"], pollCallback); err != nil {
						log.Error().Fields(map[string]interface{}{
							"action":     "send quiz poll",
							"sub action": "request resp to struct",
							"error":      err,
							"resp":       respData,
						}).Send()
					} else {
						msgId = int64(respData["message_id"].(float64))
						log.Info().Fields(map[string]interface{}{
							"action":       "send quiz poll",
							"pollCallBack": pollCallback,
						}).Send()
					}
				}

			} else {
				log.Error().Fields(map[string]interface{}{
					"action":     "send quiz poll",
					"sub action": "marshal request json",
					"error":      err,
				}).Send()
			}
		}
	}
	return pollCallback, msgId, err
}

func (b *Bot) SendVotePoll(chatId int64, question string, options []string, needAnonymous bool, multi bool) (*tgbotapi.Poll, error) {

	pollCfg := tgbotapi.NewPoll(chatId, question)
	pollCfg.AllowsMultipleAnswers = multi
	pollCfg.Options = options
	pollCfg.IsAnonymous = needAnonymous
	pollCallback := &tgbotapi.Poll{}
	resp, err := b.Request(pollCfg)
	if err != nil {
		log.Error().Fields(map[string]interface{}{
			"action": "send vote poll",
			"error":  err,
		}).Send()
	} else {
		if !resp.Ok {
			log.Error().Fields(map[string]interface{}{
				"action":       "send vote poll",
				"request resp": resp,
			}).Send()
			err = errors.New(resp.Description)
		} else {
			var respData map[string]interface{}
			var bytes []byte
			if bytes, err = resp.Result.MarshalJSON(); err == nil {
				if err = json.Unmarshal(bytes, &respData); err != nil {
					log.Error().Fields(map[string]interface{}{
						"action":     "send vote poll",
						"sub action": "unmarshal request json",
						"error":      err,
					}).Send()
				} else {
					if err = obj_util.MapToStruct(respData["poll"], pollCallback); err != nil {
						log.Error().Fields(map[string]interface{}{
							"action":     "send vote poll",
							"sub action": "request resp to struct",
							"error":      err,
							"resp":       respData,
						}).Send()
					} else {
						log.Info().Fields(map[string]interface{}{
							"action":       "send vote poll",
							"pollCallBack": pollCallback,
						}).Send()
					}
				}

			} else {
				log.Error().Fields(map[string]interface{}{
					"action":     "send vote poll",
					"sub action": "marshal request json",
					"error":      err,
				}).Send()
			}
		}
	}
	return pollCallback, err
}

func (b *Bot) StopPoll(chatId int64, msgId int64) {
	stopPoll := tgbotapi.NewStopPoll(chatId, int(msgId))
	if resp, err := b.Request(stopPoll); err != nil {
		log.Error().Fields(map[string]interface{}{
			"action": "stop poll",
			"chatId": chatId,
			"msgId":  msgId,
			"error":  err,
		}).Send()
	} else {
		if !resp.Ok {
			log.Error().Fields(map[string]interface{}{
				"action": "stop poll",
				"chatId": chatId,
				"msgId":  msgId,
				"error":  resp.Description,
			}).Send()
		}
	}
}

func (b *Bot) PinMsg(chatId int64, messageId int) error {
	res, err := b.Request(tgbotapi.PinChatMessageConfig{
		ChatID:    chatId,
		MessageID: messageId,
	})
	if err != nil {
		log.Error().Fields(map[string]interface{}{
			"action": "pin chat message",
			"name":   b.BotName,
			"error":  err,
		}).Send()
		return err
	}
	if !res.Ok {
		log.Error().Fields(map[string]interface{}{
			"action": "telegram bot kick user",
			"name":   b.BotName,
			"error":  res.Description,
		}).Send()
		err = errors.New(res.Description)
	}
	return err
}

func (b *Bot) UnPinMsg(chatId int64, messageId int) error {
	res, err := b.Request(tgbotapi.UnpinChatMessageConfig{
		ChatID:    chatId,
		MessageID: messageId,
	})
	if err != nil {
		log.Error().Fields(map[string]interface{}{
			"action": "pin chat message",
			"name":   b.BotName,
			"error":  err,
		}).Send()
		return err
	}
	if !res.Ok {
		log.Error().Fields(map[string]interface{}{
			"action": "telegram bot kick user",
			"name":   b.BotName,
			"error":  res.Description,
		}).Send()
		err = errors.New(res.Description)
	}
	return err
}

func (b *Bot) KickUser(chat *tgbotapi.Chat, user *tgbotapi.User) error {
	kick := tgbotapi.KickChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chat.ID,
			UserID: user.ID,
		},
		UntilDate:      0,
		RevokeMessages: true,
	}
	res, err := b.Request(kick)
	if err != nil {
		log.Error().Fields(map[string]interface{}{
			"action": "telegram bot kick user",
			"name":   b.BotName,
			"chat":   chat.Title,
			"user":   user.UserName,
			"error":  err,
		}).Send()
		return err
	}
	if !res.Ok {
		log.Error().Fields(map[string]interface{}{
			"action": "telegram bot kick user",
			"name":   b.BotName,
			"chat":   chat.Title,
			"user":   user.UserName,
			"error":  res.Description,
		}).Send()
		err = errors.New(res.Description)
	}
	return err
}
