package impl

import (
	"gitlab.com/rendezvousPoint/incentive_calculator/repository/database"
)

type SchedulerImpl struct {
	DbService database.DbService
}

func GetSchedulerService(dbService database.DbService) *SchedulerImpl {

	return &SchedulerImpl{
		DbService: dbService,
	}
}
