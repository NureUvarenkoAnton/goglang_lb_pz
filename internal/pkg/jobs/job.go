package jobs

import "github.com/go-co-op/gocron/v2"

type JobHandler struct {
	scheduler gocron.Scheduler
}

func NewJobHandler(scheduler gocron.Scheduler) *JobHandler {
	return &JobHandler{
		scheduler: scheduler,
	}
}
