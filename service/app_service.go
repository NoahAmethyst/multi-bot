package service

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"multi-bot/constant"
	"multi-bot/entity/entity_pb/alliance_bot_pb"
	"multi-bot/manager"
	"multi-bot/utils/log"
)

//Get discord groups
func GetGroups(req *alliance_bot_pb.GetGroupReq) *alliance_bot_pb.GetGroupResp {
	resp := &alliance_bot_pb.GetGroupResp{
		Data: &alliance_bot_pb.GetGroupRespData{},
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

	var groups []*alliance_bot_pb.Group
	resp.Data = &alliance_bot_pb.GetGroupRespData{Groups: groups}
	switch req.AppType {
	case constant.PlatTG:
		for _, singleGroup := range app.Groups {
			groups = append(groups, &alliance_bot_pb.Group{
				Id:   singleGroup.GroupId,
				Name: singleGroup.GroupName,
			})
		}

	case constant.PlatDiscord:
		if app.DiscordBotApi == nil {
			log.Error().Fields(map[string]interface{}{"action": "discord api nil", "app": app}).Send()
			//resp.Error = &alliance_bot_pb.ErrorMessage{
			//	Code:    constant.ServerErr,
			//	Message: constant.ApplicationNotFound,
			//	Detail:  fmt.Sprintf(constant.AppNotFoundDetail, req.AppId),
			//}
			return resp
		}

		guilds := app.DiscordBotApi.Session.State.Guilds
		for _, guild := range guilds {
			channelList, err := app.DiscordBotApi.Session.GuildChannels(guild.ID)
			if err != nil {
				log.Error().Fields(map[string]interface{}{"action": "discord get channel error", "error": err.Error(), "req": req, "botName": app.DiscordBotApi.Session.State.User.Username}).Send()
				return resp
			}
			for _, channel := range channelList {
				if channel.Type == discordgo.ChannelTypeGuildText {
					groups = append(groups, &alliance_bot_pb.Group{
						Id:      channel.ID,
						Name:    fmt.Sprintf("%s | %s", guild.Name, channel.Name),
						Profile: "",
					})
				}
			}
		}

		log.Info().Fields(map[string]interface{}{"action": "get groups", "groups": groups, "appId": app.Id}).Send()

		//groupList, err := applicationdb.GetTgApplicationGroups(app.Id, constant.PlatDiscord)
		//if err != nil {
		//	log.Error().Fields(map[string]interface{}{"action": "get application error", "error": err.Error()}).Send()
		//	//resp.Error = &alliance_bot_pb.ErrorMessage{
		//	//	Code:    constant.ServerErr,
		//	//	Message: constant.ApplicationNotFound,
		//	//	Detail:  constant.ServerError + err.Error(),
		//	//}
		//	//resp.Data = &alliance_bot_pb.GetGroupRespData{Groups: groups}
		//	return resp
		//}
		//log.Info().Fields(map[string]interface{}{"action": "get groups", "groups": groups, "appId": app.Id}).Send()
		//for _, g := range groupList {
		//	ch, err := app.DiscordBotApi.Session.Channel(g.GroupId)
		//	if err != nil {
		//		//log.Error().Fields(map[string]interface{}{"action": "get group error", "error": err.Error(), "chId": g.GroupId, "req": req}).Send()
		//		//resp.Error = &alliance_bot_pb.ErrorMessage{
		//		//	Code:    constant.ServerErr,
		//		//	Message: constant.ApplicationNotFound,
		//		//	Detail:  constant.ServerError + err.Error(),
		//		//}
		//		//resp.Data = &alliance_bot_pb.GetGroupRespData{Groups: groups}
		//		return resp
		//	}
		//
		//	guild, err := app.DiscordBotApi.Session.Guild(ch.GuildID)
		//	if err != nil {
		//		log.Error().Fields(map[string]interface{}{"action": "get guild error", "error": err.Error(), "guildId": ch.GuildID, "req": req}).Send()
		//		//resp.Error = &alliance_bot_pb.ErrorMessage{
		//		//	Code:    constant.ServerErr,
		//		//	Message: constant.ApplicationNotFound,
		//		//	Detail:  constant.ServerError + err.Error(),
		//		//}
		//		resp.Data = &alliance_bot_pb.GetGroupRespData{Groups: groups}
		//		return resp
		//	}
		//
		//	groups = append(groups, &alliance_bot_pb.Group{
		//		Id:   g.GroupId,
		//		Name: fmt.Sprintf("%s | %s", guild.Name, ch.Name),
		//	})
		//}

	default:

		log.Error().Msgf(constant.NoSuchBotDetail, req.AppType)
		resp.Error = &alliance_bot_pb.ErrorMessage{
			Code:    constant.ServerErr,
			Message: constant.ServerError,
			Detail:  fmt.Sprintf(constant.NoSuchBotDetail, req.AppType),
		}
	}

	resp.Data = &alliance_bot_pb.GetGroupRespData{Groups: groups}

	return resp
}
