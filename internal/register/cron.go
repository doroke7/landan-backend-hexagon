package register

import (
	"log"

	"github.com/robfig/cron/v3"

	"example/internal/container"
)

func CronInit(oContainer *container.CronContainer) *cron.Cron {
	oCron := cron.New()

	if _, err := oCron.AddFunc("* * * * *", oContainer.CronAppUser.IncreaseBalance); err != nil {
		log.Fatalf("cron: failed to register CronAppUser.IncreaseBalance job: %v", err)
	}

	return oCron
}
