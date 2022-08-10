package service

import (
	"fmt"
	"multi-bot/app_context"
	"multi-bot/callback"
	"multi-bot/constant"
	"multi-bot/entity/entity_pb/alliance_bot_pb"
	"multi-bot/entity/entity_pb/tg_quiz_pb"
	"multi-bot/manager"
	"multi-bot/memory"
	"multi-bot/utils/encrypt_util"
	"multi-bot/utils/log"
	"strconv"
	"time"
)

func CreateQuiz(req *alliance_bot_pb.CreateQuizReq) *alliance_bot_pb.CreateQuizResp {
	log.Info().Fields(map[string]interface{}{
		"action": "create quiz start",
		"req":    req,
	}).Send()

	resp := &alliance_bot_pb.CreateQuizResp{
		Data: &alliance_bot_pb.CreateQuizRespData{},
	}
	app, ok := manager.GetApplication(req.AppId)
	if !ok {
		log.Error().Msgf("application:%s not found", req.AppId)
		resp.Error = &alliance_bot_pb.ErrorMessage{
			Code:    constant.ServerErr,
			Message: constant.ApplicationNotFound,
			Detail:  fmt.Sprintf(constant.AppNotFoundDetail, req.AppId),
		}
		return resp
	}

	if len(req.Quiz) == 0 {
		resp.Error = &alliance_bot_pb.ErrorMessage{
			Code:    constant.ParamsErr,
			Message: constant.ParamsError,
			Detail:  "quiz can't be empty",
		}
		return resp
	}

	quizId := req.QuizId
	if len(quizId) == 0 {
		quizId = encrypt_util.GenerateUuid(true)
	}

	var group string

	if len(req.GroupId) > 0 {
		group = req.GroupId
	} else {
		group = app.Groups[0].GroupId
	}

	switch req.AppType {
	case constant.PlatTG:
		tgGroupId, _ := strconv.ParseInt(group, 10, 64)
		if app.TgBotApi == nil {

			log.Error().Msgf(constant.AppBotNotSetDetail, app.Name)
			resp.Error = &alliance_bot_pb.ErrorMessage{
				Code:    constant.ServerErr,
				Message: constant.ServerError,
				Detail:  fmt.Sprintf(constant.AppBotNotSetDetail, app.Name),
			}
			return resp
		}
		var quizIds []string
		for _, question := range req.Quiz {
			if sentQuestion, msgId, err := app.TgBotApi.SendQuizPoll(tgGroupId, question.Question, question.Options, question.CorrectIndex, req.NeedAnonymous); err != nil {
				log.Error().Fields(map[string]interface{}{
					"action": "send question",
					"error":  err,
				}).Send()
				resp.Error = &alliance_bot_pb.ErrorMessage{
					Code:    constant.ServerErr,
					Message: constant.ServerError,
					Detail:  err.Error(),
				}
				return resp
			} else {
				memory.PollMemory.SetQuestionMsgIdMap(sentQuestion.ID, msgId)
				memory.PollMemory.SetQuestionChatMap(sentQuestion.ID, tgGroupId)
				memory.PollMemory.SetPollMemory(sentQuestion.ID, sentQuestion)
				memory.PollMemory.SetQuestionQuizMap(sentQuestion.ID, quizId)
				quizIds = append(quizIds, sentQuestion.ID)
			}
			time.Sleep(500 * time.Millisecond)
		}
		memory.PollMemory.SetQuizQuestionsMap(quizId, quizIds)

		if req.ActiveTime > 0 {
			go func() {
				memory.PollMemory.AddQuizJob(quizId, &tg_quiz_pb.Job{
					AppId:       app.Id,
					QuizId:      quizId,
					AppType:     constant.PlatTG,
					ExecuteTime: time.Now().Unix() + req.ActiveTime,
				})
				time.Sleep(time.Duration(req.ActiveTime) * time.Second)
				stopQuiz(app.AppContext, quizId, constant.PlatTG)
			}()
		}
	case constant.PlatDiscord:
		//not support

	default:
		log.Error().Msgf(constant.NoSuchBotDetail, req.AppType)
		resp.Error = &alliance_bot_pb.ErrorMessage{
			Code:    constant.ServerErr,
			Message: constant.ServerError,
			Detail:  fmt.Sprintf(constant.NoSuchBotDetail, req.AppType),
		}
		return resp
	}

	resp.Data = &alliance_bot_pb.CreateQuizRespData{
		QuizId:  quizId,
		GroupId: group,
	}

	return resp
}

func StopQuiz(req *alliance_bot_pb.StopQuizReq) *alliance_bot_pb.CommonResp {
	log.Info().Fields(map[string]interface{}{
		"action": "stop quiz start",
		"req":    req,
	}).Send()
	resp := &alliance_bot_pb.CommonResp{}

	app, ok := manager.GetApplication(req.AppId)
	if !ok {
		resp.Error = &alliance_bot_pb.ErrorMessage{
			Code:    constant.ServerErr,
			Message: constant.ApplicationNotFound,
			Detail:  fmt.Sprintf(constant.AppBotNotSetDetail, req.AppId),
		}
		return resp
	}

	switch req.AppType {
	case constant.PlatTG:
		if app.TgBotApi == nil {
			log.Error().Msgf(constant.AppBotNotSetDetail, app.Name)
			resp.Error = &alliance_bot_pb.ErrorMessage{
				Code:    constant.ServerErr,
				Message: constant.ServerError,
				Detail:  fmt.Sprintf(constant.AppBotNotSetDetail, app.Name),
			}
			return resp
		}
		stopQuiz(app.AppContext, req.QuizId, constant.PlatTG)
	case constant.PlatDiscord:
		//not support
	default:
		log.Error().Msgf(constant.NoSuchBotDetail, req.AppType)
		resp.Error = &alliance_bot_pb.ErrorMessage{
			Code:    constant.ServerErr,
			Message: constant.ServerError,
			Detail:  fmt.Sprintf(constant.NoSuchBotDetail, req.AppType),
		}
		return resp
	}

	resp.Success = true
	return resp
}

func stopQuiz(appCtx *app_context.AppContext, quizId string, appType int32) {
	questionIds := memory.PollMemory.GetQuizQuestionIds(quizId)
	for _, questionId := range questionIds {
		stopQuestion(appCtx, questionId)
	}
	callback.QuizStop(appCtx, quizId, appType)
	memory.PollMemory.ClearQuiz(quizId)
	log.Info().Fields(map[string]interface{}{
		"action": "stop question",
		"quizId": quizId,
		"app":    appCtx.Name,
	}).Send()
}

func stopQuestion(appCtx *app_context.AppContext, quizId string) {
	chatId := memory.PollMemory.GetQuestionChatId(quizId)
	msgId := memory.PollMemory.GetQuestionMsgId(quizId)
	appCtx.TgBotApi.StopPoll(chatId, msgId)
	memory.PollMemory.ClearQuestionMemory(quizId)
}
