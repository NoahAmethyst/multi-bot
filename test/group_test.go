package test

import (
	"github.com/tristan-club/wizard/pkg/util"
	"multi-bot/entity/entity_pb/alliance_bot_pb"
	"multi-bot/service"
	"testing"
)

func TestGetGroups(t *testing.T) {
	resp := service.GetGroups(&alliance_bot_pb.GetGroupReq{
		AppId:   "discord-quiz-test",
		AppType: 2,
	})
	t.Log(util.FastMarshal(resp))
}
