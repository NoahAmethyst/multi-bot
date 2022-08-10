package cron

import (
	"github.com/cenkalti/backoff"
	"github.com/robfig/cron"
	"multi-bot/utils/log"
	"reflect"
	"time"
)

const (
	// 初始重试时间
	initialInterval = time.Second * 3
	// 最大重试次数
	maxRetryTime = 3
)

const (
	JobDurationHourly = "@hourly"
	JobDurationDaily  = "@daily"
)

// AddCronJob 时任务执行器
func AddCronJob(cronJob func() error, jobDuration string, preHandle bool) error {

	if preHandle {
		if err := cronJob(); err != nil {
			log.Error().Msgf("cronjob job %s err %s", reflect.TypeOf(cronJob).Name(), err.Error())
			return err
		}
	}

	c := cron.New()
	err := c.AddFunc(jobDuration, func() {
		if err := cronJob(); err != nil {
			log.Error().Msgf("cronjob job %s error %s", reflect.TypeOf(cronJob), err.Error())
		} else {
			log.Info().Msgf("cronjob job %s success", reflect.TypeOf(cronJob))
		}
	})
	if err != nil {
		log.Error().Msgf("add cronjob job error %s, check your input ", err.Error())
		return err
	}

	c.Start()
	return nil
}

// AddCronJobWithBackoff 带退避重试算法的定时任务执行器
func AddCronJobWithBackoff(cronJob func() error, jobDuration string) error {

	if err := cronJob(); err != nil {
		log.Error().Msgf("cronjob job %s err %s", reflect.TypeOf(cronJob()).Name(), err.Error())
		return err
	}

	c := cron.New()
	err := c.AddFunc(jobDuration, func() {
		backOffHandler(cronJob)
	})
	if err != nil {
		return err
	}

	c.Start()
	return nil
}

// backOffHandler 退避算法执行失败任务
func backOffHandler(cronJob func() error) {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = initialInterval
	err := backoff.Retry(cronJob, backoff.WithMaxRetries(b, maxRetryTime))
	if err != nil {
		log.Error().Msgf("cronjob job %s  backoff handler after retry %d times, got error %s", reflect.TypeOf(cronJob), maxRetryTime, err.Error())
	} else {
		log.Info().Msgf("cronjob job %s backoff handler success", reflect.TypeOf(cronJob))
	}
}
