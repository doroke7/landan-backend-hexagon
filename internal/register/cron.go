package register

import (
	"github.com/robfig/cron/v3"

	"example/internal/container"
)

func CronInit(oContainer *container.CronContainer) *cron.Cron {
	oCron := cron.New()

	return oCron
}
