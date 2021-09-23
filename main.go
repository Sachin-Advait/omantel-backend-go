package main

import (
	"gitlab.com/rendezvousPoint/incentive_calculator/commons"
	"gitlab.com/rendezvousPoint/incentive_calculator/config"
	"gitlab.com/rendezvousPoint/incentive_calculator/core/crons"
	"gitlab.com/rendezvousPoint/incentive_calculator/core/handler"
	"gitlab.com/rendezvousPoint/incentive_calculator/infra"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils"
)

const (
	baseURL = "/sales/incentive"
	health  = "/health"
)

func init() {
	utils.SetLoggerObj(commons.LOG_FILE_NAME)
}

func main() {

	cronService := crons.NewCronServices(&config.GetInstance().Database)
	cronContext := crons.NewCronContext(cronService)

	go cronContext.StartCrons()

	serverContext := infra.NewServerContext(&config.GetInstance().Server)

	trxService := handler.NewSalesIncentiveService(&config.GetInstance().Database)
	trxContext := handler.NewSalesIncentiveContext(trxService)

	defer trxContext.SafeClose()

	serverContext.Mount(baseURL, trxContext.NewSalesIncentiveRouter())
	serverContext.Mount(health, handler.NewHealthCheckRouter())

	// Start server
	address := config.GetInstance().Server.GetAddressString()
	serverContext.StartServer(address)
}
