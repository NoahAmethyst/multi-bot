package service

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"multi-bot/constant"
	"multi-bot/entity/entity_pb/alliance_bot_pb"
	"multi-bot/manager"
	"multi-bot/telegram/keyboard"
	"multi-bot/utils/log"
	"strconv"
	"time"
)

var proxy map[alliance_bot_pb.MsgType]func(req *alliance_bot_pb.SendMsgReq) *alliance_bot_pb.SendMsgResp

func SendMsgProxy(req *alliance_bot_pb.SendMsgReq) *alliance_bot_pb.SendMsgResp {

	handler, ok := proxy[req.MsgType]
	if !ok {
		handler = SendMsg
	}
	return handler(req)
}

func SendMsg(req *alliance_bot_pb.SendMsgReq) *alliance_bot_pb.SendMsgResp {
	resp, app, group, done, ikm := preHandle(req)
	if done {
		return resp
	}

	content := req.Content

	switch req.AppType {
	case constant.PlatTG:
		tgGroupId, _ := strconv.ParseInt(group, 10, 64)
		if len(req.MentionUserIds) > 0 {
			var mention string
			for _, mentionId := range req.MentionUserIds {
				if res, err := app.TgBotApi.GetChatMember(tgbotapi.GetChatMemberConfig{
					ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
						ChatID: tgGroupId,
						UserID: mentionId,
					},
				}); err == nil && res.User != nil {
					mention += fmt.Sprintf("@%s ", res.User.UserName)
				}
			}
			content = fmt.Sprintf(content, mention)
		}

		if thisMsg, err := app.TgBotApi.SendMsg(&tgbotapi.Chat{
			ID: tgGroupId,
		}, content, ikm, req.Markdown, req.Preview); err != nil {
			resp.Error = &alliance_bot_pb.ErrorMessage{
				Code:    constant.ServerErr,
				Message: constant.ServerError,
				Detail:  err.Error(),
			}

		} else {
			app.TgBotApi.DelDeadMsg(&thisMsg, time.Duration(req.Deadline)*time.Second)
			resp.Data = &alliance_bot_pb.SendMsgRespData{
				GroupId:   strconv.FormatInt(tgGroupId, 10),
				MessageId: int64(thisMsg.MessageID),
			}
		}
	case constant.PlatDiscord:
		if len(req.MentionUserIds) > 0 {
			var mention string
			for _, mentionId := range req.MentionUserIds {
				mention += fmt.Sprintf("<@%d>", mentionId)
			}
			content = fmt.Sprintf(content, mention)
		}

		title := req.Title
		if title == "" {
			title = ""
		}

		embed := &discordgo.MessageEmbed{
			Type:        "rich",
			Title:       title,
			Description: content,
		}

		messageSend := &discordgo.MessageSend{
			Content:         "",
			TTS:             false,
			Components:      nil,
			Files:           nil,
			AllowedMentions: nil,
			Reference:       nil,
			File:            nil,
			Embed:           embed,
		}

		if bt, ok := ikm.(*discordgo.Button); ok {
			messageSend.Components = []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						bt,
					},
				},
			}
		}

		msg, err := app.DiscordBotApi.Session.ChannelMessageSendComplex(group, messageSend)
		if err != nil {
			resp.Error = &alliance_bot_pb.ErrorMessage{
				Code:    constant.ServerErr,
				Message: constant.ServerError,
				Detail:  err.Error(),
			}
			log.Error().Fields(map[string]interface{}{"action": "discord send msg error", "error": err.Error(), "req": req}).Send()
		} else {
			messageId, err := strconv.ParseInt(msg.ID, 10, 64)
			if err != nil {
				messageId = 0
				log.Error().Fields(map[string]interface{}{"action": "invalid message id", "error": err.Error(), "msgId": msg.ID, "req": req}).Send()
			}
			resp.Data = &alliance_bot_pb.SendMsgRespData{
				GroupId:   group,
				MessageId: messageId,
			}
		}
	}
	return resp
}

func SendPhoto(req *alliance_bot_pb.SendMsgReq) *alliance_bot_pb.SendMsgResp {
	resp, app, group, done, ikm := preHandle(req)
	if done {
		return resp
	}

	switch req.AppType {
	case constant.PlatTG:
		tgGroupId, _ := strconv.ParseInt(group, 10, 64)
		if thisMsg, err := app.TgBotApi.SendPhoto(&tgbotapi.Chat{
			ID: tgGroupId,
		}, req.PhotoUrl, ikm); err != nil {
			resp.Error = &alliance_bot_pb.ErrorMessage{
				Code:    constant.ServerErr,
				Message: constant.ServerError,
				Detail:  err.Error(),
			}

		} else {
			app.TgBotApi.DelDeadMsg(&thisMsg, time.Duration(req.Deadline)*time.Second)
			resp.Data = &alliance_bot_pb.SendMsgRespData{
				GroupId:   group,
				MessageId: int64(thisMsg.MessageID),
			}
		}

	case constant.PlatDiscord:
		title := req.Title
		if title == "" {
			title = "Note"
		}

		embed := &discordgo.MessageEmbed{
			Type:        "rich",
			Title:       title,
			Description: req.Content,
			Image: &discordgo.MessageEmbedImage{
				URL:      req.PhotoUrl,
				ProxyURL: "",
				Width:    0,
				Height:   0,
			},
		}

		messageSend := &discordgo.MessageSend{
			Content:         "",
			TTS:             false,
			Components:      nil,
			Files:           nil,
			AllowedMentions: nil,
			Reference:       nil,
			File:            nil,
			Embed:           embed,
		}

		if bt, ok := ikm.(*discordgo.Button); ok {
			messageSend.Components = []discordgo.MessageComponent{
				bt,
			}
		}

		msg, err := app.DiscordBotApi.Session.ChannelMessageSendComplex(group, messageSend)
		if err != nil {
			resp.Error = &alliance_bot_pb.ErrorMessage{
				Code:    constant.ServerErr,
				Message: constant.ServerError,
				Detail:  err.Error(),
			}
		} else {
			messageId, err := strconv.ParseInt(msg.ID, 10, 64)
			if err != nil {
				messageId = 0
				log.Error().Fields(map[string]interface{}{"action": "invalid message id", "error": err.Error(), "msgId": msg.ID, "req": req}).Send()
			}
			resp.Data = &alliance_bot_pb.SendMsgRespData{
				GroupId:   group,
				MessageId: messageId,
			}

		}
	}

	return resp
}

//Todo test
//Get discord users
func GetUsers(req *alliance_bot_pb.GetUsersReq) *alliance_bot_pb.GetUsersResp {
	resp := &alliance_bot_pb.GetUsersResp{
		Data: &alliance_bot_pb.GetUserRespData{Users: nil},
	}
	app, ok := manager.GetApplication(req.AppId)
	if !ok {
		resp.Error = &alliance_bot_pb.ErrorMessage{
			Code:    constant.ServerErr,
			Message: constant.ApplicationNotFound,
			Detail:  fmt.Sprintf(constant.AppNotFoundDetail, req.AppId),
		}
		return resp
	}
	data := map[int64]*alliance_bot_pb.TelegramUser{}
	switch req.AppType {
	case constant.PlatTG:
		if app.TgBotApi == nil {
			log.Error().Msgf(constant.AppBotNotSetDetail, app.Name)
			resp.Error = &alliance_bot_pb.ErrorMessage{
				Code:    constant.ServerErr,
				Message: constant.ServerError,
				Detail:  fmt.Sprintf(constant.AppBotNotSetDetail, req.AppId),
			}
			return resp
		}

		if len(req.UserId) == 0 {
			log.Warn().Msgf("req userId is empty")
		}
		for _, userId := range req.UserId {
			chatId, _ := strconv.ParseInt(req.GroupId, 10, 64)
			member, err := app.TgBotApi.GetChatMember(tgbotapi.GetChatMemberConfig{
				ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
					ChatID: chatId,
					UserID: userId,
				},
			})
			if err != nil {
				log.Error().Msgf("get chat member error:%s", err)
				continue
			}
			data[userId] = &alliance_bot_pb.TelegramUser{Name: member.User.UserName}
		}

		if len(data) == 0 {
			log.Warn().Msgf("user data is empty")
		}

		resp.Data.Users = data
	case constant.PlatDiscord:
		err := func(req *alliance_bot_pb.GetUsersReq) error {
			if app.DiscordBotApi == nil {
				return fmt.Errorf(constant.AppBotNotSetDetail, req.AppId)
			}
			channel, err := app.DiscordBotApi.Session.Channel(req.GroupId)
			if err != nil {
				log.Error().Fields(map[string]interface{}{"action": "get channel error", "error": err.Error(), "req": req}).Send()
				return fmt.Errorf("get channel info error %s", err.Error())
			}
			for _, v := range req.UserId {
				member, err := app.DiscordBotApi.Session.GuildMember(channel.GuildID, strconv.FormatInt(v, 10))
				if err != nil {
					log.Error().Fields(map[string]interface{}{"action": "get user error", "error": err.Error(), "userId": v}).Send()
					return fmt.Errorf("get user(id:%d) info error", v)
				}
				if member.Nick != "" {
					data[v] = &alliance_bot_pb.TelegramUser{Name: member.Nick}
				} else {
					data[v] = &alliance_bot_pb.TelegramUser{Name: member.User.Username}
				}
			}
			return nil
		}(req)
		if err != nil {
			log.Error().Fields(map[string]interface{}{"action": "get members error", "error": err.Error(), "req": req}).Send()
			resp.Error = &alliance_bot_pb.ErrorMessage{
				Code:    constant.ServerErr,
				Message: constant.ServerError,
				Detail:  err.Error(),
			}
			return resp
		} else {
			resp.Data.Users = data
		}
	default:
		log.Error().Msgf(constant.NoSuchBotDetail, req.AppType)
		resp.Error = &alliance_bot_pb.ErrorMessage{
			Code:    constant.ServerErr,
			Message: constant.ServerError,
			Detail:  fmt.Sprintf(constant.NoSuchBotDetail, req.AppType),
		}
		return resp

	}

	return resp
}

func GetBotInfo(req *alliance_bot_pb.GetBotInfoReq) *alliance_bot_pb.GetBotInfoResp {
	resp := &alliance_bot_pb.GetBotInfoResp{
		Data: &alliance_bot_pb.GetBotInfoRespData{},
	}

	app, ok := manager.GetApplication(req.AppId)
	if !ok {
		resp.Error = &alliance_bot_pb.ErrorMessage{
			Code:    constant.ServerErr,
			Message: constant.ApplicationNotFound,
			Detail:  fmt.Sprintf(constant.AppNotFoundDetail, req.AppId),
		}
		return resp
	}

	var hasBot bool
	resp.Data = &alliance_bot_pb.GetBotInfoRespData{
		TgBotInfo:      nil,
		DiscordBotInfo: nil,
	}
	if app.AppContext.TgBotApi != nil {
		hasBot = true
		resp.Data.TgBotInfo = &alliance_bot_pb.TelegramBotInfo{
			BotName: app.TgBotApi.BotName,
			Type:    constant.PlatTG,
		}

		resp.Data.Type = constant.PlatTG
		resp.Data.BotName = app.TgBotApi.BotName
	}

	if app.AppContext.DiscordBotApi != nil {
		hasBot = true
		resp.Data.DiscordBotInfo = &alliance_bot_pb.DiscordBotInfo{
			BotName: app.DiscordBotApi.GetName(),
			Type:    constant.PlatDiscord,
			BotId:   app.DiscordBotApi.Session.State.User.ID,
		}
	}

	if !hasBot {
		log.Error().Msgf(constant.AppBotNotSetDetail, app.Name)
		resp.Error = &alliance_bot_pb.ErrorMessage{
			Code:    constant.ServerErr,
			Message: constant.ServerError,
			Detail:  fmt.Sprintf(constant.AppBotNotSetDetail, req.AppId),
		}
	}

	return resp
}

func preHandle(req *alliance_bot_pb.SendMsgReq) (*alliance_bot_pb.SendMsgResp, manager.Application, string, bool, interface{}) {
	log.Info().Fields(map[string]interface{}{
		"action": "send message",
		"req":    req,
	}).Send()
	resp := &alliance_bot_pb.SendMsgResp{
		Data: &alliance_bot_pb.SendMsgRespData{},
	}

	app, ok := manager.GetApplication(req.AppId)
	if !ok {
		resp.Error = &alliance_bot_pb.ErrorMessage{
			Code:    constant.ServerErr,
			Message: constant.ApplicationNotFound,
			Detail:  fmt.Sprintf(constant.AppNotFoundDetail, req.AppId),
		}
		return resp, manager.Application{}, "", true, nil
	}
	var ikm interface{}
	var group string
	switch req.AppType {
	case constant.PlatTG:
		if app.AppContext.TgBotApi == nil {
			log.Error().Msgf(constant.AppBotNotSetDetail, app.Name)
			resp.Error = &alliance_bot_pb.ErrorMessage{
				Code:    constant.ServerErr,
				Message: constant.ServerError,
				Detail:  fmt.Sprintf(constant.AppBotNotSetDetail, req.AppId),
			}
			return resp, manager.Application{}, "", true, nil
		}

		if len(req.GroupId) > 0 {
			group = req.GroupId
		} else {
			group = app.Groups[0].GroupId
		}

		switch req.InlineMarkup {
		case alliance_bot_pb.InlineMarkupType_None:
			break
		case alliance_bot_pb.InlineMarkupType_ForwardBot:
			ikm, _ = keyboard.NewForwardPrivateKeyBoard(constant.Participate, app.TgBotApi)
		default:
			break
		}

	case constant.PlatDiscord:
		if app.AppContext.DiscordBotApi == nil {
			log.Error().Msgf(constant.AppBotNotSetDetail, app.Name)
			resp.Error = &alliance_bot_pb.ErrorMessage{
				Code:    constant.ServerErr,
				Message: constant.ServerError,
				Detail:  fmt.Sprintf(constant.AppBotNotSetDetail, req.AppId),
			}
			return resp, manager.Application{}, "", true, nil
		}

		if req.InlineMarkup == alliance_bot_pb.InlineMarkupType_ForwardBot {
			ikm = &discordgo.Button{
				Label:    constant.ButtonParticipate,
				Style:    discordgo.PrimaryButton,
				Disabled: false,
				//Emoji:    discordgo.ComponentEmoji{},
				//URL:      fmt.Sprintf("https://t.me/%s", app.DiscordBotApi.Token),
				CustomID: constant.DiscordCustomIdStart,
			}
		}

		group = req.GroupId
	default:
		log.Error().Msgf(constant.NoSuchBotDetail, req.AppType)
		resp.Error = &alliance_bot_pb.ErrorMessage{
			Code:    constant.ServerErr,
			Message: constant.ServerError,
			Detail:  fmt.Sprintf(constant.NoSuchBotDetail, req.AppType),
		}

		return resp, manager.Application{}, "", true, nil
	}

	return resp, app, group, false, ikm
}

func init() {
	proxy = map[alliance_bot_pb.MsgType]func(req *alliance_bot_pb.SendMsgReq) *alliance_bot_pb.SendMsgResp{
		alliance_bot_pb.MsgType_Content: SendMsg,
		alliance_bot_pb.MsgType_Photo:   SendPhoto,
	}
}
