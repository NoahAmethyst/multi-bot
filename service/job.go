package service

import (
	"multi-bot/entity/entity_pb/alliance_bot_pb"
	"multi-bot/entity/entity_pb/tg_quiz_pb"
	"multi-bot/memory"
	"multi-bot/utils/log"
	"time"
)

func InitJob() {
	for k, v := range memory.PollMemory.GetJobs() {
		log.Info().Fields(map[string]interface{}{
			"action": "execute unhandled job",
			"job":    v,
		}).Send()
		go func(id string, job *tg_quiz_pb.Job) {
			now := time.Now().Unix()
			if job.ExecuteTime <= now {
				StopQuiz(&alliance_bot_pb.StopQuizReq{
					AppId:   job.AppId,
					QuizId:  id,
					AppType: job.AppType,
				})
			} else {
				time.Sleep(time.Duration(job.ExecuteTime-now) * time.Second)
				StopQuiz(&alliance_bot_pb.StopQuizReq{
					AppId:   job.AppId,
					QuizId:  id,
					AppType: job.AppType,
				})
			}
		}(k, v)
	}
}
