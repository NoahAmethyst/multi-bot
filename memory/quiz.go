package memory

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"multi-bot/entity/entity_pb/tg_quiz_pb"
	"sync"
)

type pollRecord struct {
	quizQuestionMemory    *tg_quiz_pb.QuizQuestionMemory
	quizParticipateMemory *tg_quiz_pb.QuizParticipateMemory
	quizPlatMemory        *tg_quiz_pb.QuizPlatMemory
	jobs                  *tg_quiz_pb.Jobs
	sync.RWMutex
}

const (
	maxWaitTime = 60
)

var PollMemory pollRecord

func (p *pollRecord) SetQuizQuestionsMap(quizId string, questionIds []string) {
	p.Lock()
	defer p.Unlock()
	//p.quizQuestionsMap[quizId] = questionIds
	p.quizQuestionMemory.QuizQuestion[quizId] = &tg_quiz_pb.Questions{Ids: questionIds}
}

func (p *pollRecord) SetQuestionMsgIdMap(questionId string, msgId int64) {
	p.Lock()
	defer p.Unlock()
	//p.questionMsgMap[questionId] = msgId
	p.quizPlatMemory.QuestionMsg[questionId] = msgId
}

func (p *pollRecord) SetQuestionChatMap(questionId string, chatId int64) {
	p.Lock()
	defer p.Unlock()
	//p.chatMemory[questionId] = chatId
	p.quizPlatMemory.ChatMemory[questionId] = chatId
}

func (p *pollRecord) GetQuestionChatId(questionId string) int64 {
	p.RLock()
	defer p.RUnlock()
	//return p.chatMemory[questionId]
	return p.quizPlatMemory.ChatMemory[questionId]
}

func (p *pollRecord) GetQuestionMsgId(quizId string) int64 {
	p.RLock()
	defer p.RUnlock()
	//return p.questionMsgMap[quizId]
	return p.quizPlatMemory.QuestionMsg[quizId]
}

func (p *pollRecord) GetQuizQuestionIds(questionId string) []string {
	p.RLock()
	defer p.RUnlock()
	//return p.quizQuestionsMap[questionId]
	return p.quizQuestionMemory.QuizQuestion[questionId].Ids
}

func (p *pollRecord) SetPollMemory(pollId string, poll *tgbotapi.Poll) {
	p.Lock()
	defer p.Unlock()
	if _, ok := p.quizPlatMemory.Memory[pollId]; ok {
		return
	}
	p.quizPlatMemory.Memory[pollId] = &tg_quiz_pb.Poll{
		Id:            poll.ID,
		CorrectOption: int64(poll.CorrectOptionID),
		PollType:      poll.Type,
	}
}

func (p *pollRecord) GetPollMemory(pollId string) *tg_quiz_pb.Poll {
	p.RLock()
	defer p.RUnlock()
	//return p.memory[pollId]
	return p.quizPlatMemory.Memory[pollId]
}

func (p *pollRecord) SetQuestionCorrectUsers(quizId string, uids ...int64) {
	p.Lock()
	defer p.Unlock()
	if v, ok := p.quizParticipateMemory.QuestionCorrectUser[quizId]; !ok {
		p.quizParticipateMemory.QuestionCorrectUser[quizId] = &tg_quiz_pb.QuestionCorrectUsers{Ids: uids}
	} else {
		v.Ids = append(v.Ids, uids...)
		p.quizParticipateMemory.QuestionCorrectUser[quizId] = v
	}
}

func (p *pollRecord) GetCorrectUser(quizId string) []int64 {
	p.RLock()
	defer p.RUnlock()
	//return p.questionCorrectMemory[quizId]
	if v, ok := p.quizParticipateMemory.QuestionCorrectUser[quizId]; !ok {
		return nil
	} else {
		return v.Ids
	}

}

func (p *pollRecord) SetQuestionQuizMap(quizId string, questionId string) {
	p.Lock()
	defer p.Unlock()
	//p.questionQuizMap[quizId] = questionId
	p.quizQuestionMemory.QuestionQuiz[quizId] = questionId
}

func (p *pollRecord) GetQuestionQuizId(quizId string) (string, bool) {
	p.RLock()
	defer p.RUnlock()
	//v, ok := p.questionQuizMap[quizId]
	v, ok := p.quizQuestionMemory.QuestionQuiz[quizId]
	return v, ok
}

func (p *pollRecord) SetQuizCorrect(questionId string, userId int64, num int) {
	p.Lock()
	defer p.Unlock()
	if v, ok := p.quizParticipateMemory.QuizCorrectUser[questionId]; !ok {
		p.quizParticipateMemory.QuizCorrectUser[questionId] = &tg_quiz_pb.QuizUserCorrect{CorrectUser: map[int64]int64{userId: int64(num)}}
	} else {
		v.CorrectUser[userId] = int64(num)
	}
}

func (p *pollRecord) GetQuizCorrect(questionId string, userId int64) int {
	p.RLock()
	defer p.RUnlock()
	var correctNum int
	if v, ok := p.quizParticipateMemory.QuizCorrectUser[questionId]; !ok {
		correctNum = 0
	} else {
		correctNum = int(v.CorrectUser[userId])
	}
	return correctNum
}

func (p *pollRecord) GetQuizParticipateNum(quizId string) int {
	p.RLock()
	defer p.RUnlock()
	//return int64(len(p.quizParticipate[quizId]))
	if v, ok := p.quizParticipateMemory.QuizParticipate[quizId]; !ok {
		return 0
	} else {
		return len(v.Participate)
	}
}

func (p *pollRecord) SetQuizParticipate(quizId string, userId int64) {
	p.Lock()
	defer p.Unlock()
	if v, ok := p.quizParticipateMemory.QuizParticipate[quizId]; !ok {
		p.quizParticipateMemory.QuizParticipate[quizId] = &tg_quiz_pb.Participate{Participate: map[int64]bool{userId: true}}
	} else {
		v.Participate[userId] = true
	}
}

func (p *pollRecord) AddQuizJob(quizId string, job *tg_quiz_pb.Job) {
	p.Lock()
	defer p.Unlock()
	p.jobs.Jobs[quizId] = job
}

func (p *pollRecord) GetJobs() map[string]*tg_quiz_pb.Job {
	p.RLock()
	defer p.RUnlock()
	jobs := map[string]*tg_quiz_pb.Job{}
	for k, v := range p.jobs.Jobs {
		jobs[k] = v
	}
	return jobs
}

func (p *pollRecord) ClearQuestionMemory(questionId string) {
	p.Lock()
	defer p.Unlock()
	delete(p.quizQuestionMemory.QuestionQuiz, questionId)
	delete(p.quizPlatMemory.QuestionMsg, questionId)
	delete(p.quizPlatMemory.ChatMemory, questionId)
	delete(p.quizPlatMemory.Memory, questionId)
	delete(p.quizParticipateMemory.QuestionCorrectUser, questionId)
}

func (p *pollRecord) ClearQuiz(quizId string) {
	p.Lock()
	defer p.Unlock()
	delete(p.quizQuestionMemory.QuizQuestion, quizId)
	delete(p.quizParticipateMemory.QuizCorrectUser, quizId)
	delete(p.quizParticipateMemory.QuizParticipate, quizId)
	delete(p.jobs.Jobs, quizId)
}

func (p *pollRecord) SaveQuizQuestion(v *tg_quiz_pb.QuizQuestionMemory) {
	p.Lock()
	defer p.Unlock()
	for k, subV := range v.QuizQuestion {
		p.quizQuestionMemory.QuizQuestion[k] = subV
	}
	for k, subV := range v.QuestionQuiz {
		p.quizQuestionMemory.QuestionQuiz[k] = subV
	}
}

func (p *pollRecord) SaveQuizParticipate(v *tg_quiz_pb.QuizParticipateMemory) {
	p.Lock()
	defer p.Unlock()
	for k, subV := range v.QuizParticipate {
		p.quizParticipateMemory.QuizParticipate[k] = subV
	}

	for k, subV := range v.QuizCorrectUser {
		p.quizParticipateMemory.QuizCorrectUser[k] = subV
	}

	for k, subV := range v.QuestionCorrectUser {
		p.quizParticipateMemory.QuestionCorrectUser[k] = subV
	}
}

func (p *pollRecord) SaveQuizPlat(v *tg_quiz_pb.QuizPlatMemory) {
	p.Lock()
	defer p.Unlock()
	for k, subV := range v.Memory {
		p.quizPlatMemory.Memory[k] = subV
	}

	for k, subV := range v.ChatMemory {
		p.quizPlatMemory.ChatMemory[k] = subV
	}

	for k, subV := range v.QuestionMsg {
		p.quizPlatMemory.QuestionMsg[k] = subV
	}
}

func (p *pollRecord) SaveQuizJobs(v *tg_quiz_pb.Jobs) {
	p.Lock()
	defer p.Unlock()

	for k, subV := range v.Jobs {
		p.jobs.Jobs[k] = subV
	}
}

func InitMemory() {
	//code init memory from persistent data source

}

func (p *pollRecord) SaveMemory() error {
	//code here save memory to persistent data source
	return nil
}

func init() {

	PollMemory = pollRecord{
		quizQuestionMemory: &tg_quiz_pb.QuizQuestionMemory{
			QuizQuestion: map[string]*tg_quiz_pb.Questions{},
			QuestionQuiz: map[string]string{},
		},
		quizParticipateMemory: &tg_quiz_pb.QuizParticipateMemory{
			QuizCorrectUser:     map[string]*tg_quiz_pb.QuizUserCorrect{},
			QuestionCorrectUser: map[string]*tg_quiz_pb.QuestionCorrectUsers{},
			QuizParticipate:     map[string]*tg_quiz_pb.Participate{},
		},
		quizPlatMemory: &tg_quiz_pb.QuizPlatMemory{
			QuestionMsg: map[string]int64{},
			ChatMemory:  map[string]int64{},
			Memory:      map[string]*tg_quiz_pb.Poll{},
		},
		jobs:    &tg_quiz_pb.Jobs{Jobs: map[string]*tg_quiz_pb.Job{}},
		RWMutex: sync.RWMutex{},
	}
}
