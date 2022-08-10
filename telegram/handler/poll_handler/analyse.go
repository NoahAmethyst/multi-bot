package poll_handler

import (
	"fmt"
	"multi-bot/app_context"
	"multi-bot/callback"
	"multi-bot/memory"
	tg_api "multi-bot/telegram/api"
	"multi-bot/utils/log"

	"multi-bot/entity/entity_pb/tg_quiz_pb"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Analyse(appCtx *app_context.AppContext, pollAnswer *tgbotapi.PollAnswer) {
	if poll := memory.PollMemory.GetPollMemory(pollAnswer.PollID); poll == nil {
		log.Error().Fields(map[string]interface{}{
			"action": "get poll memory",
			"error":  "not found",
		}).Send()
		return
	} else {
		switch poll.PollType {
		case tg_api.PollQuiz:
			AnalyseQuiz(appCtx, pollAnswer)
		case tg_api.PollRegular:
			AnalyseVote(appCtx, pollAnswer)
		default:
			log.Error().Fields(map[string]interface{}{
				"action": "analyse poll answer",
				"error":  fmt.Sprintf("no such type:%s", poll.PollType),
			}).Send()
		}
	}
}

func AnalyseQuiz(appCtx *app_context.AppContext, pollAnswer *tgbotapi.PollAnswer) {
	poll := memory.PollMemory.GetPollMemory(pollAnswer.PollID)
	if poll == nil {
		log.Error().Fields(map[string]interface{}{
			"action": "get poll memory",
			"error":  "not found",
		}).Send()
		return
	}

	isCorrect := CheckQuizCorrect(poll, pollAnswer)
	if isCorrect {
		memory.PollMemory.SetQuestionCorrectUsers(pollAnswer.PollID, pollAnswer.User.ID)
	}
	callback.QuizCorrectUser(appCtx, poll.Id, &pollAnswer.User, isCorrect)
}

func AnalyseVote(appCtx *app_context.AppContext, pollAnswer *tgbotapi.PollAnswer) {

}

func CheckQuizCorrect(poll *tg_quiz_pb.Poll, pollAnswer *tgbotapi.PollAnswer) bool {
	for _, v := range pollAnswer.OptionIDs {
		if int(poll.CorrectOption) == v {
			return true
		}
	}
	return false
}
