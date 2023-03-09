package scheduler

import (
	"gitlab.com/rendezvousPoint/incentive_calculator/business/scheduler"
)

var scheduleProvide scheduler.Scheduler

func InitScheduler(provider scheduler.Scheduler) {
	scheduleProvide = provider
}
