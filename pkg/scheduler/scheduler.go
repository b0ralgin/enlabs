package scheduler

import (
	"enlabs/pkg/account"
	"time"
)

type scheduler struct {
	am account.Corrector
}

func RunScheduler(am account.Corrector, t time.Duration) error {

}
