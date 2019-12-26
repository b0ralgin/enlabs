package scheduler

import (
	"enlabs/pkg/account"
	"errors"

	"github.com/jasonlvhit/gocron"
)

//RunScheduler run scheduler
func RunScheduler(am account.Corrector, t uint64) error {
	scheduler := gocron.NewScheduler()
	scheduler.Every(t).Minutes().Do(am.CorrectBalance)
	<-scheduler.Start()
	return errors.New("scheduler has stopped")
}
