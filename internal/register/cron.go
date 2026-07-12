package register

import (
	"log"

	"github.com/robfig/cron/v3"

	"example/internal/container"
)

func CronInit(oContainer *container.Container) *cron.Cron {
	oCron := cron.New()

	if _, err := oCron.AddFunc("* * * * *", oContainer.CronUser.AddUser); err != nil {
		log.Fatalf("cron: failed to register CronUser.AddUser job: %v", err)
	}

	return oCron
}
