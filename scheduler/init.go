package scheduler

import (
	"github.com/robfig/cron/v3"
	"gitlab.com/rendezvousPoint/incentive_calculator/adapters/scheduler"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/configs"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
)

func StartCronJob() error {

	applicationConfig, err := configs.Get(constant.ApplicationConfigFile)
	if err != nil {
		return err
	}
	masterCronTime := applicationConfig.GetString(constant.MasterCronTime)

	scheduler.PreLoadIntrinsicValues()
	scheduler.PrintIntrinsicValues()

	cronJob := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

	_, err = cronJob.AddFunc(masterCronTime, scheduler.LoadIntrinsicValues)
	if err != nil {
		return err
	}

	cronJob.Start()

	return nil
}
