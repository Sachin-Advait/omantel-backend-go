package impl

import (
	"gitlab.com/rendezvousPoint/incentive_calculator/repository/database"
	databaseClient "gitlab.com/rendezvousPoint/incentive_calculator/utils/database"
)

type DbServiceImpl struct {
	DBClient databaseClient.Database
}

func GetDBServiceInstance(dbClient databaseClient.Database) database.DbService {
	return &DbServiceImpl{
		DBClient: dbClient,
	}
}
