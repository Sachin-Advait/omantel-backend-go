package models

import (
	"errors"
	"fmt"
	"strings"

	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func ValidateUpsertRegion(region []*CreateRegionRequest) error {

	if len(region) == 0 || region[0] == nil {
		return errors.New("invalid request")
	}

	for _, eachRecord := range region {

		if eachRecord.RegionCode == "" {
			return errors.New("invalid regionCode")
		}
		eachRecord.RegionCode = strings.ToUpper(strings.TrimSpace(eachRecord.RegionCode))

		if eachRecord.RegionName == "" {
			return errors.New("invalid regionName")
		}
		eachRecord.RegionName = strings.ToUpper(strings.TrimSpace(eachRecord.RegionName))

		if eachRecord.Quarter < 1 || eachRecord.Quarter > 4 {
			return errors.New("invalid quarter")
		}

		if eachRecord.Year < 2022 {
			return errors.New("invalid year")
		}
	}

	return nil
}

func ValidateUpsertOutlet(outlet []*CreateOutletRequest) error {

	if len(outlet) == 0 || outlet[0] == nil {
		return errors.New("invalid request")
	}

	for _, eachRecord := range outlet {

		if eachRecord.OutletCode == "" {
			return errors.New("invalid outletCode")
		}
		eachRecord.OutletCode = strings.ToUpper(strings.TrimSpace(eachRecord.OutletCode))

		if eachRecord.OutletName == "" {
			return errors.New("invalid outletName")
		}
		eachRecord.OutletName = strings.ToUpper(strings.TrimSpace(eachRecord.OutletName))

		if eachRecord.RegionCode == "" {
			return errors.New("invalid regionCode")
		}
		eachRecord.RegionCode = strings.ToUpper(strings.TrimSpace(eachRecord.RegionCode))

		if eachRecord.RegionName == "" {
			return errors.New("invalid regionName")
		}
		eachRecord.RegionName = strings.ToUpper(strings.TrimSpace(eachRecord.RegionName))

		if eachRecord.Lat == 0 {
			return errors.New("invalid latitude")
		}

		if eachRecord.Lang == 0 {
			return errors.New("invalid langitude")
		}

		if eachRecord.Quarter < 1 || eachRecord.Quarter > 4 {
			return errors.New("invalid quarter")
		}

		if eachRecord.Year < 2022 {
			return errors.New("invalid year")
		}
	}

	return nil
}

func ValidateUpsertProduct(product []*CreateProductRequest) error {

	if len(product) == 0 || product[0] == nil {
		return errors.New("invalid request")
	}

	for _, eachRecord := range product {

		if eachRecord.ProductCode == "" {
			return errors.New("invalid productCode")
		}
		eachRecord.ProductCode = strings.ToUpper(strings.TrimSpace(eachRecord.ProductCode))

		if eachRecord.ProductName == "" {
			return errors.New("invalid productName")
		}
		eachRecord.ProductName = strings.ToUpper(strings.TrimSpace(eachRecord.ProductName))

		if eachRecord.Quarter < 1 || eachRecord.Quarter > 4 {
			return errors.New("invalid quarter")
		}

		if eachRecord.Year < 2022 {
			return errors.New("invalid year")
		}
	}

	return nil
}

func ValidateUpsertPlan(plan []*CreatePlanRequest) error {

	if len(plan) == 0 || plan[0] == nil {
		return errors.New("invalid request")
	}

	for _, eachRecord := range plan {

		if eachRecord.ProductCode == "" {
			return errors.New("invalid productCode")
		}
		eachRecord.ProductCode = strings.ToUpper(strings.TrimSpace(eachRecord.ProductCode))

		if eachRecord.PlanCode == "" {
			return errors.New("invalid planCode")
		}
		eachRecord.PlanCode = strings.ToUpper(strings.TrimSpace(eachRecord.PlanCode))

		if eachRecord.PlanName == "" {
			return errors.New("invalid planName")
		}
		eachRecord.PlanName = strings.ToUpper(strings.TrimSpace(eachRecord.PlanName))

		if eachRecord.PlanPrice == 0 {
			return errors.New("invalid planPrice")
		}

		if eachRecord.SIPPlanRate == 0 {
			return errors.New("invalid sipRate")
		}

		if eachRecord.Quarter < 1 || eachRecord.Quarter > 4 {
			return errors.New("invalid quarter")
		}

		if eachRecord.Year < 2022 {
			return errors.New("invalid year")
		}
	}

	return nil
}

func ValidateUpsertTarget(target []*CreateTargetRequest) error {

	if len(target) == 0 || target[0] == nil {
		return errors.New("invalid request")
	}

	for _, eachRecord := range target {

		if eachRecord.TargetCode == "" {
			return errors.New("invalid targetCode")
		}
		eachRecord.TargetCode = strings.ToUpper(strings.TrimSpace(eachRecord.TargetCode))

		if eachRecord.RegionName == "" {
			return errors.New("invalid regionName")
		}
		eachRecord.RegionName = strings.ToUpper(strings.TrimSpace(eachRecord.RegionName))

		if eachRecord.OutletName == "" {
			return errors.New("invalid outleteName")
		}
		eachRecord.OutletName = strings.ToUpper(strings.TrimSpace(eachRecord.OutletName))

		if eachRecord.HBB == 0 {
			return errors.New("invalid hbbTarget")
		}

		if eachRecord.Hayyak == 0 {
			return errors.New("invalid hayyakTarget")
		}

		if eachRecord.Baqati == 0 {
			return errors.New("invalid baqatiTarget")
		}

		if eachRecord.Upgrade == 0 {
			return errors.New("invalid upgradeTarget")
		}

		if eachRecord.Month < 1 || eachRecord.Month > 12 {
			return errors.New("invalid month")
		}

		quarter := objecthelper.GetQuarterFromMonth(eachRecord.Month)
		if eachRecord.Quarter != 0 && eachRecord.Quarter != quarter {
			return errors.New("invalid quarter")
		}
		eachRecord.Quarter = quarter

		if eachRecord.Year < 2022 {
			return errors.New("invalid year")
		}
	}

	return nil
}

func ValidateUpsertVOCCriteria(voc []*CreateVOCCriteriaRequest) error {

	if len(voc) == 0 || voc[0] == nil {
		return errors.New("invalid request")
	}

	for _, eachRecord := range voc {

		if eachRecord.Unsatisfactory == 0 {
			return errors.New("invalid unsatisfactory")
		}

		if eachRecord.AlmostAchieved == 0 {
			return errors.New("invalid almostAchieved")
		}

		if eachRecord.Achieved == 0 {
			return errors.New("invalid achieved")
		}

		if eachRecord.WellAchieved == 0 {
			return errors.New("invalid wellAchieved")
		}

		if eachRecord.Outstanding == 0 {
			return errors.New("invalid outstanding")
		}

		if eachRecord.Year < 2022 {
			return errors.New("invalid year")
		}
	}

	return nil
}

func ValidateUpsertVOCMaxPay(voc []*CreateVOCMaxPayRequest) error {

	if len(voc) == 0 || voc[0] == nil {
		return errors.New("invalid request")
	}

	for _, eachRecord := range voc {

		if _, ok := constant.EmployeePositionMap[eachRecord.Position]; !ok {
			return errors.New("invalid position-" + eachRecord.Position)
		}

		if eachRecord.MaxPay == 0 {
			return errors.New("invalid maxPay")
		}

		if eachRecord.Year < 2022 {
			return errors.New("invalid year")
		}
	}

	return nil
}

func ValidateUpsertVOCUser(voc []*CreateVOCUserRequest) error {

	if len(voc) == 0 || voc[0] == nil {
		return errors.New("invalid request")
	}

	for _, eachRecord := range voc {

		if eachRecord.UserId == "" {
			return errors.New("invalid userId")
		}

		if eachRecord.AchievedPercent == 0 {
			return errors.New("invalid achievedPercent")
		}

		if eachRecord.Year < 2022 {
			return errors.New("invalid year")
		}
	}

	return nil
}

func ValidateFetchRegion(region *FetchDeleteRegionRequest) error {

	if region.Quarter < 1 || region.Quarter > 4 {
		return errors.New("invalid quarter")
	}

	if region.Year <= 0 {
		return errors.New("invalid year")
	}

	return nil
}

func ValidateFetchOutlet(region *FetchDeleteOutletRequest) error {

	if region.Quarter < 1 || region.Quarter > 4 {
		return errors.New("invalid quarter")
	}

	if region.Year <= 0 {
		return errors.New("invalid year")
	}

	return nil
}

func ValidateFetchProduct(region *FetchDeleteProductRequest) error {

	if region.Quarter < 1 || region.Quarter > 4 {
		return errors.New("invalid quarter")
	}

	if region.Year <= 0 {
		return errors.New("invalid year")
	}

	return nil
}

func ValidateFetchPlan(region *FetchDeletePlanRequest) error {

	if region.Quarter < 1 || region.Quarter > 4 {
		return errors.New("invalid quarter")
	}

	if region.Year <= 0 {
		return errors.New("invalid year")
	}

	return nil
}

func ValidateFetchTarget(target *FetchDeleteTargetRequest) error {

	if target.Year <= 0 {
		return errors.New("invalid year")
	}

	if len(target.TargetCodes) > 0 && len(target.OutletCode) > 0 {
		return errors.New("enter either target code or outlet, not both")
	}

	if len(target.TargetCodes) > 0 {
		for i, v := range target.TargetCodes {
			if v == "" {
				return errors.New("invalid target code")
			}
			target.TargetCodes[i] = strings.ToUpper(strings.TrimSpace(v))
		}
	}

	if len(target.OutletCode) > 0 {
		for i, v := range target.OutletCode {
			if v == "" {
				return errors.New("invalid outlet code")
			}
			target.OutletCode[i] = strings.ToUpper(strings.TrimSpace(v))
		}
	}

	if target.Month == 0 && (target.Quarter >= 1 || target.Quarter <= 4) {
		return nil
	}

	if target.Quarter == 0 && (target.Month >= 1 || target.Month <= 12) {
		target.Quarter = objecthelper.GetQuarterFromMonth(target.Month)
		return nil
	}

	if (target.Quarter >= 1 || target.Quarter <= 4) && (target.Month >= 1 || target.Month <= 12) {
		quarter := objecthelper.GetQuarterFromMonth(target.Month)
		if target.Quarter == quarter {
			return nil
		}
	}

	return errors.New("invalid month or quarter")
}

func ValidateFetchVOCCriteria(voc *FetchDeleteVOCCriteriaRequest) error {

	if voc == nil {
		return errors.New("invalid request")
	}

	if voc.Year < 2022 {
		return errors.New("invalid year")
	}

	return nil
}

func ValidateFetchVOCUser(voc *FetchDeleteVOCUserRequest) error {

	if voc == nil {
		return errors.New("invalid request")
	}

	if len(voc.UserIds) == 0 {
		return errors.New("zero userIds")
	}

	if voc.Year < 2022 {
		return errors.New("invalid year")
	}

	for index, user := range voc.UserIds {
		if user == "" {
			return errors.New("blank userId at index-" + fmt.Sprintf("%v", index))
		}
	}

	return nil
}

func ValidateDeleteRegion(region *FetchDeleteRegionRequest) error {

	if len(region.RegionCodes) == 0 {
		return errors.New("invalid regionCodes")
	}

	for i, v := range region.RegionCodes {
		if v == "" {
			return errors.New("empty regionCode")
		}
		region.RegionCodes[i] = strings.ToUpper(strings.TrimSpace(v))
	}

	if region.Quarter < 1 || region.Quarter > 4 {
		return errors.New("invalid quarter")
	}

	if region.Year <= 0 {
		return errors.New("invalid year")
	}

	return nil
}

func ValidateDeleteOutlet(region *FetchDeleteOutletRequest) error {

	if len(region.OutletCodes) == 0 {
		return errors.New("invalid outletCodes")
	}

	for i, v := range region.OutletCodes {
		if v == "" {
			return errors.New("empty outletCode")
		}
		region.OutletCodes[i] = strings.ToUpper(strings.TrimSpace(v))
	}

	if region.Quarter < 1 || region.Quarter > 4 {
		return errors.New("invalid quarter")
	}

	if region.Year <= 0 {
		return errors.New("invalid year")
	}

	return nil
}

func ValidateDeleteProduct(region *FetchDeleteProductRequest) error {

	if len(region.ProductCodes) == 0 {
		return errors.New("invalid productCodes")
	}

	for i, v := range region.ProductCodes {
		if v == "" {
			return errors.New("empty productCode")
		}
		region.ProductCodes[i] = strings.ToUpper(strings.TrimSpace(v))
	}

	if region.Quarter < 1 || region.Quarter > 4 {
		return errors.New("invalid quarter")
	}

	if region.Year <= 0 {
		return errors.New("invalid year")
	}

	return nil
}

func ValidateDeletePlan(region *FetchDeletePlanRequest) error {

	if len(region.PlanCodes) == 0 {
		return errors.New("invalid planCodes")
	}

	for i, v := range region.PlanCodes {
		if v == "" {
			return errors.New("empty planCode")
		}
		region.PlanCodes[i] = strings.ToUpper(strings.TrimSpace(v))
	}

	if region.Quarter < 1 || region.Quarter > 4 {
		return errors.New("invalid quarter")
	}

	if region.Year <= 0 {
		return errors.New("invalid year")
	}

	return nil
}

func ValidateDeleteTarget(target *FetchDeleteTargetRequest) error {

	if target.Year <= 0 {
		return errors.New("invalid year")
	}

	if len(target.TargetCodes) > 0 && len(target.OutletCode) > 0 {
		return errors.New("enter either target code or outlet, not both")
	}

	if len(target.TargetCodes) > 0 {
		for i, v := range target.TargetCodes {
			if v == "" {
				return errors.New("invalid target code")
			}
			target.TargetCodes[i] = strings.ToUpper(strings.TrimSpace(v))
		}
	}

	if len(target.OutletCode) > 0 {
		for i, v := range target.OutletCode {
			if v == "" {
				return errors.New("invalid outlet")
			}
			target.OutletCode[i] = strings.ToUpper(strings.TrimSpace(v))
		}
	}

	if target.Month == 0 && (target.Quarter >= 1 || target.Quarter <= 4) {
		return nil
	}

	if target.Quarter == 0 && (target.Month >= 1 || target.Month <= 12) {
		target.Quarter = objecthelper.GetQuarterFromMonth(target.Month)
		return nil
	}

	if (target.Quarter >= 1 || target.Quarter <= 4) && (target.Month >= 1 || target.Month <= 12) {
		quarter := objecthelper.GetQuarterFromMonth(target.Month)
		if target.Quarter == quarter {
			return nil
		}
	}

	return errors.New("invalid month or quarter")
}

func ValidateDeleteVOCCriteria(voc *FetchDeleteVOCCriteriaRequest) error {

	if voc == nil {
		return errors.New("invalid request")
	}

	if voc.Year < 2022 {
		return errors.New("invalid year")
	}

	return nil
}

func ValidateDeleteVOCMaxPay(voc []*FetchDeleteVOCMaxPayRequest) error {

	if len(voc) == 0 {
		return errors.New("invalid request")
	}

	for i, val := range voc {
		if _, ok := constant.EmployeePositionMap[val.Position]; !ok {
			return fmt.Errorf("invalid position, index-%v, position-%v", i, val.Position)
		}

		if val.Year < 2022 {
			return fmt.Errorf("invalid year, index-%v, year-%v", i, val.Year)
		}
	}

	return nil
}

func ValidateDeleteVOCUser(voc *FetchDeleteVOCUserRequest) error {

	if voc == nil {
		return errors.New("invalid request")
	}

	if len(voc.UserIds) == 0 {
		return errors.New("zero userIds")
	}

	if voc.Year < 2022 {
		return errors.New("invalid year")
	}

	for index, user := range voc.UserIds {
		if user == "" {
			return errors.New("blank userId at index-" + fmt.Sprintf("%v", index))
		}
	}

	return nil
}
