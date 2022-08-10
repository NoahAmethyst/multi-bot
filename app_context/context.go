package app_context

import (
	"google.golang.org/grpc"
	"multi-bot/constant"
	"multi-bot/db/applicationdb"
	discord_api "multi-bot/discord/api"
	tg_api "multi-bot/telegram/api"
	"strconv"
)

type AppContext struct {
	Id            string
	Name          string
	Url           string
	Conn          *grpc.ClientConn
	TgBotApi      *tg_api.Bot
	DiscordBotApi *discord_api.Bot
	Groups        []*applicationdb.TgApplicationGroup
}

func (a *AppContext) SetApplicationBot(tgBot *tg_api.Bot) {
	a.TgBotApi = tgBot
}

func (a *AppContext) CheckInGroup(groupId int64) bool {
	if len(a.Groups) == 0 {
		return false
	}
	for _, thisGroup := range a.Groups {
		if thisGroup.GroupId == strconv.FormatInt(groupId, 10) {
			return true
		}
	}
	return false
}

func (a *AppContext) AddGroup(groupId int64, groupName string) {
	newGroup := &applicationdb.TgApplicationGroup{
		AppId:     a.Id,
		GroupId:   strconv.FormatInt(groupId, 10),
		GroupName: groupName,
		IsDelete:  false,
		BotType:   constant.PlatTG,
	}
	a.Groups = append(a.Groups, newGroup)

	_ = applicationdb.InsertTgApplicationGroup(newGroup)
}

func (a *AppContext) ChangeGroup(groupId int64, newName string) {
	var newGroup *applicationdb.TgApplicationGroup
	for i, group := range a.Groups {
		if group.GroupId == strconv.FormatInt(groupId, 10) {
			newGroup = group
			newGroup.GroupName = newName
			a.Groups[i] = newGroup
			break
		}
	}

	if newGroup != nil {
		_ = applicationdb.InsertTgApplicationGroup(newGroup)
	}
}

func (a *AppContext) DelGroup(groupId int64) {
	var index int
	var delGroup *applicationdb.TgApplicationGroup
	var exist bool
	for i, group := range a.Groups {
		if group.GroupId == strconv.FormatInt(groupId, 10) {
			delGroup = group
			index = i
			exist = true
			break
		}
	}
	if exist {
		a.Groups = append(a.Groups[:index], a.Groups[index+1:]...)
		_ = applicationdb.DelTgApplicationGroup(delGroup)
	}
}
