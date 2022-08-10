package callback

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"multi-bot/app_context"
	"multi-bot/constant"
	"multi-bot/entity/entity_pb/alliance_bot_pb"
	"multi-bot/memory"
	"multi-bot/utils/http_util"
	"multi-bot/utils/log"
	"strconv"
	"time"
)

const (
	maxRetryTime = 6
)

//Correct user callback
func QuizCorrectUser(appCtx *app_context.AppContext, questionId string, user *tgbotapi.User, correct bool) {

	quizId, ok := memory.PollMemory.GetQuestionQuizId(questionId)
	if !ok {
		log.Warn().Msgf("quiz - question not found")
		return
	}

	memory.PollMemory.SetQuizParticipate(quizId, user.ID)

	if !correct {
		return
	}

	quizIds := memory.PollMemory.GetQuizQuestionIds(quizId)
	totalQuizNum := len(quizIds)
	correctNum := memory.PollMemory.GetQuizCorrect(quizId, user.ID)
	if correctNum+1 == totalQuizNum {
		log.Info().Msgf("user %s correct all quizs", user.UserName)
		if appCtx.Conn != nil {
			//Todo grpc interface{}

		} else {
			resp := map[string]interface{}{}
			data := &alliance_bot_pb.QuizCallbackReq{
				QuizId:       quizId,
				OpenId:       user.ID,
				AppId:        appCtx.Id,
				AppType:      constant.PlatTG,
				CallbackType: alliance_bot_pb.QuizCallbackType_CorrectUser,
			}
			if err := http_util.PostJSON(appCtx.Url, data, nil, &resp); err != nil {
				retryPostJson(appCtx, data, err, resp)
			} else {
				log.Info().Fields(map[string]interface{}{
					"action":    "callback handle quiz:correct user",
					"app":       appCtx.Name,
					"appSvcUri": appCtx.Url,
					"data":      data,
					"resp":      resp,
				}).Send()
			}
		}
	} else {
		memory.PollMemory.SetQuizCorrect(quizId, user.ID, correctNum+1)
	}
}

//Stop quiz callback
func QuizStop(appCtx *app_context.AppContext, quizId string, appType int32) {
	resp := map[string]interface{}{}
	data := &alliance_bot_pb.QuizCallbackReq{
		QuizId:         quizId,
		AppId:          appCtx.Id,
		CallbackType:   alliance_bot_pb.QuizCallbackType_Stop,
		ParticipateNum: int64(memory.PollMemory.GetQuizParticipateNum(quizId)),
		AppType:        appType,
	}

	if err := http_util.PostJSON(appCtx.Url, data, nil, &resp); err != nil {
		retryPostJson(appCtx, data, err, resp)
	} else {
		log.Info().Fields(map[string]interface{}{
			"action":    "callback handle quiz:stop",
			"app":       appCtx.Name,
			"appSvcUri": appCtx.Url,
			"data":      data,
			"resp":      resp,
		}).Send()
	}
}

func QuizGenerate(appCtx *app_context.AppContext, groupId int64, appType int32) {
	resp := map[string]interface{}{}
	data := &alliance_bot_pb.QuizCallbackReq{
		CallbackType: alliance_bot_pb.QuizCallbackType_Generate,
		GroupId:      strconv.FormatInt(groupId, 10),
		AppId:        appCtx.Id,
		AppType:      appType,
	}

	if err := http_util.PostJSON(appCtx.Url, data, nil, &resp); err != nil {
		retryPostJson(appCtx, data, err, resp)
	} else {
		log.Info().Fields(map[string]interface{}{
			"action":    "callback handle quiz:generate",
			"app":       appCtx.Name,
			"appSvcUri": appCtx.Url,
			"data":      data,
			"resp":      resp,
		}).Send()
	}
}

func retryPostJson(appCtx *app_context.AppContext, data *alliance_bot_pb.QuizCallbackReq, err error, resp map[string]interface{}) {
	go func() {
		for retryTime := 0; retryTime < maxRetryTime; retryTime++ {
			time.Sleep(time.Duration(retryTime+1) * 10 * time.Second)
			log.Warn().Fields(map[string]interface{}{
				"action":    "callback handle quiz",
				"app":       appCtx.Name,
				"appSvcUri": appCtx.Url,
				"data":      data,
				"retryTime": retryTime + 1,
			}).Send()
			if err = http_util.PostJSON(appCtx.Url, data, nil, &resp); err == nil {
				return
			}
		}
		log.Error().Fields(map[string]interface{}{
			"action":    "callback handle quiz",
			"app":       appCtx.Name,
			"appSvcUri": appCtx.Url,
			"data":      data,
			"error":     "retry max time failed",
		}).Send()
	}()
}
