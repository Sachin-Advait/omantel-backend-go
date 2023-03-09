package user

import "gitlab.com/rendezvousPoint/incentive_calculator/business/user"

var userAPIProvider user.UserProvider

func InitUserAPI(provider user.UserProvider) {
	userAPIProvider = provider
}
