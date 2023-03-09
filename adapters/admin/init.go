package admin

import (
	"gitlab.com/rendezvousPoint/incentive_calculator/business/admin"
)

var adminAPIProvider admin.AdminProvider

func InitAdminAPI(provider admin.AdminProvider) {
	adminAPIProvider = provider
}
