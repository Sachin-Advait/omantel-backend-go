package main

import (
	"encoding/json"
	"fmt"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
)

func main() {

	Parsing()

	now := time.Now()
	time.Sleep(2 * time.Second)

	fmt.Println("sub", time.Now().Sub(now))

	for {
		time.Sleep(1 * time.Second)
		fmt.Println("application is running - " + time.Now().String())
	}
}

func Parsing() {
	hbbSaleMap := make(map[string]int, 0)
	str := `[{"current_plan":"WIRELESS UNLIMITED DATA - 8 MBPS - 25 RO","qty":24},{"current_plan":"WIRELESS 20MBPS - 25 OMR","qty":8},{"current_plan":"WIRELESS UNLIMITED DATA - 4 MBPS - 20 RO","qty":7}]`

	fmt.Println("str", str)

	saleCount := make([]*models.SaleCountPerPlan, 0)
	if err := json.Unmarshal([]byte(str), &saleCount); err == nil {
		for _, sale := range saleCount {
			hbbSaleMap[sale.CurrentPlan] += sale.Quantity
		}

		saleCount = saleCount[:0]
	} else {
		fmt.Printf("agent err: %+v\n", err)
	}

	fmt.Println("map", hbbSaleMap)
}
