package bot_service

import (
	"context"
	"multi-bot/constant"
	"multi-bot/entity/entity_pb/alliance_bot_pb"
	"multi-bot/manager"
	"multi-bot/service"
	"multi-bot/utils/log"
)

type Server struct{}

func (s Server) CreateQuiz(_ context.Context, req *alliance_bot_pb.CreateQuizReq) (*alliance_bot_pb.CreateQuizResp, error) {
	if req.AppType == 0 {
		req.AppType = constant.PlatTG
	}
	return service.CreateQuiz(req), nil
}

func (s Server) HelloWorld(_ context.Context, req *alliance_bot_pb.HelloWorldReq) (resp *alliance_bot_pb.CommonResp, err error) {

	resp = &alliance_bot_pb.CommonResp{
		Error:   nil,
		Success: true,
	}

	return resp, err
}

func (s Server) StopQuiz(_ context.Context, req *alliance_bot_pb.StopQuizReq) (*alliance_bot_pb.CommonResp, error) {
	if req.AppType == 0 {
		req.AppType = constant.PlatTG
	}
	return service.StopQuiz(req), nil
}

func (s Server) GetGroups(_ context.Context, req *alliance_bot_pb.GetGroupReq) (*alliance_bot_pb.GetGroupResp, error) {
	if req.AppType == 0 {
		req.AppType = constant.PlatTG
	}

	return service.GetGroups(req), nil
}

func (s Server) SendMsg(_ context.Context, req *alliance_bot_pb.SendMsgReq) (*alliance_bot_pb.SendMsgResp, error) {
	if req.AppType == 0 {
		req.AppType = constant.PlatTG
	}
	return service.SendMsgProxy(req), nil
}

func (s Server) GetBotInfo(_ context.Context, req *alliance_bot_pb.GetBotInfoReq) (*alliance_bot_pb.GetBotInfoResp, error) {
	return service.GetBotInfo(req), nil
}

func (s Server) GetUsers(_ context.Context, req *alliance_bot_pb.GetUsersReq) (*alliance_bot_pb.GetUsersResp, error) {
	if req.AppType == 0 {
		req.AppType = constant.PlatTG
	}
	return service.GetUsers(req), nil
}

func (s Server) RestartBot(_ context.Context, _ *alliance_bot_pb.RestartReq) (*alliance_bot_pb.CommonResp, error) {
	log.Info().Msg("restart all bots")
	manager.StartAllBots()
	return &alliance_bot_pb.CommonResp{
		Error:   nil,
		Success: true,
	}, nil
}
