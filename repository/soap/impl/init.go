package impl

import (
	"gitlab.com/rendezvousPoint/incentive_calculator/repository/soap"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/http"
)

type SOAPServiceImpl struct {
	HttpClient http.HttpClient
}

func GetSOAPServiceInstance(httpClient http.HttpClient) soap.SOAPService {
	return &SOAPServiceImpl{
		HttpClient: httpClient,
	}
}
