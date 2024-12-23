package jobs

import (
	"context"
	"time"

	"github.com/go-co-op/gocron/v2"
)

const MONTH = time.Hour * 24 * 31

type iDelteUserService interface {
	DelteMarkedUsers(ctx context.Context, t time.Time)
}

func (h *JobHandler) RegisterClearUsers(userService iDelteUserService) {
	h.scheduler.NewJob(
		gocron.CronJob("0 3 * * *", false),
		gocron.NewTask(func(userService iDelteUserService) {
			// delete users that have been in delete state for 2 moths
			userService.DelteMarkedUsers(context.Background(), time.Now().Add(-2*MONTH))
		}, userService),
	)
}
