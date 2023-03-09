package health

import "gitlab.com/rendezvousPoint/incentive_calculator/business/health"

var healthProvider health.HealthCheckProvider

func InitHealthCheck(provider health.HealthCheckProvider) {
	healthProvider = provider
}
