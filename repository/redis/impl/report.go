package impl

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *CacheServiceImpl) InsertUserReport(isQuarterReport bool, month, quarter, year int, userReport *models.UserReportDB, keyTTL time.Duration) *errors.Error {

	if userReport == nil || userReport.UserId == "" {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "insertUserReport: cache missing details",
		}
	}

	bytes, err := objecthelper.MarshalJSON(userReport)
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "insertUserReport: cache marshal failed - " + err.Error(),
		}
	}

	var key string

	if isQuarterReport {
		key = fmt.Sprintf(constant.RedisUserReportKey, userReport.UserId, 0, quarter, year)
	} else {
		key = fmt.Sprintf(constant.RedisUserReportKey, userReport.UserId, month, 0, year)
	}

	err = impl.RedisClient.Set(context.Background(), key, string(bytes), keyTTL).Err()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "insertUserReport: cache operation failed - " + err.Error(),
		}
	}

	return nil
}

func (impl *CacheServiceImpl) GetUserReport(userId string, isQuarterReport bool, month, quarter, year int) (*models.UserReportDB, *errors.Error) {

	if userId == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "getUserReport: cache missing details",
		}
	}

	var key string

	if isQuarterReport {
		key = fmt.Sprintf(constant.RedisUserReportKey, userId, 0, quarter, year)
	} else {
		key = fmt.Sprintf(constant.RedisUserReportKey, userId, month, 0, year)
	}

	data, err := impl.RedisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "getUserReport: cache operation failed - " + err.Error(),
		}
	}

	if err == redis.Nil || data == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusNotFound,
			Code:       errors.UserNotFoundCode,
			Message:    errors.UserNotFound,
			Details:    "getUserReport: cache no record found",
		}
	}

	report := &models.UserReportDB{}

	err = objecthelper.UnmarshalJSON([]byte(data), &report)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "getUserReport: cache unmarshal failed - " + err.Error(),
		}
	}

	return report, nil
}

func (impl *CacheServiceImpl) DelUserReport(userId string, isQuarterReport bool, month, quarter, year int) *errors.Error {

	if userId == "" {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "delUserReport: cache missing details",
		}
	}

	var key string

	if isQuarterReport {
		key = fmt.Sprintf(constant.RedisUserReportKey, userId, 0, quarter, year)
	} else {
		key = fmt.Sprintf(constant.RedisUserReportKey, userId, month, 0, year)
	}

	err := impl.RedisClient.Del(context.Background(), key).Err()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "delUserReport: cache operation failed - " + err.Error(),
		}
	}

	return nil
}

func (impl *CacheServiceImpl) InsertOutletReport(isQuarterReport bool, month, quarter, year int, outletReport *models.OutletReportDB, keyTTL time.Duration) *errors.Error {

	if outletReport == nil || outletReport.Outlet == "" || outletReport.Region == "" {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "insertOutletReport: cache missing details",
		}
	}

	bytes, err := objecthelper.MarshalJSON(outletReport)
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "insertOutletReport: cache marshal failed - " + err.Error(),
		}
	}

	var key string

	if isQuarterReport {
		key = fmt.Sprintf(constant.RedisOutletRegionReportKey, outletReport.Outlet, outletReport.Region, 0, quarter, year)
	} else {
		key = fmt.Sprintf(constant.RedisOutletRegionReportKey, outletReport.Outlet, outletReport.Region, month, 0, year)
	}

	err = impl.RedisClient.Set(context.Background(), key, string(bytes), keyTTL).Err()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "insertOutletReport: cache operation failed - " + err.Error(),
		}
	}

	return nil
}

func (impl *CacheServiceImpl) GetOutletReport(outlet, region string, isQuarterReport bool, month, quarter, year int) (*models.OutletReportDB, *errors.Error) {

	if outlet == "" || region == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "getOutletReport: cache missing details",
		}
	}

	var key string

	if isQuarterReport {
		key = fmt.Sprintf(constant.RedisOutletRegionReportKey, outlet, region, 0, quarter, year)
	} else {
		key = fmt.Sprintf(constant.RedisOutletRegionReportKey, outlet, region, month, 0, year)
	}

	data, err := impl.RedisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "getOutletReport: cache operation failed - " + err.Error(),
		}
	}

	if err == redis.Nil || data == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusNotFound,
			Code:       errors.UserNotFoundCode,
			Message:    errors.UserNotFound,
			Details:    "getOutletReport: cache no record found",
		}
	}

	report := &models.OutletReportDB{}

	err = objecthelper.UnmarshalJSON([]byte(data), &report)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "getOutletReport: cache unmarshal failed - " + err.Error(),
		}
	}

	return report, nil
}

func (impl *CacheServiceImpl) DelOutletReport(outlet, region string, isQuarterReport bool, month, quarter, year int) *errors.Error {

	if outlet == "" || region == "" {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "delOutletReport: cache missing details",
		}
	}

	var key string

	if isQuarterReport {
		key = fmt.Sprintf(constant.RedisOutletRegionReportKey, outlet, region, 0, quarter, year)
	} else {
		key = fmt.Sprintf(constant.RedisOutletRegionReportKey, outlet, region, month, 0, year)
	}

	err := impl.RedisClient.Del(context.Background(), key).Err()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "delOutletReport: cache operation failed - " + err.Error(),
		}
	}

	return nil
}

func (impl *CacheServiceImpl) InsertRegionReport(isQuarterReport bool, month, quarter, year int, regionReport *models.RegionReportDB, keyTTL time.Duration) *errors.Error {

	if regionReport == nil || regionReport.Region == "" {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "insertRegionReport: cache missing details",
		}
	}

	bytes, err := objecthelper.MarshalJSON(regionReport)
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "insertRegionReport: cache marshal failed - " + err.Error(),
		}
	}

	var key string

	if isQuarterReport {
		key = fmt.Sprintf(constant.RedisOutletRegionReportKey, "", regionReport.Region, 0, quarter, year)
	} else {
		key = fmt.Sprintf(constant.RedisOutletRegionReportKey, "", regionReport.Region, month, 0, year)
	}

	err = impl.RedisClient.Set(context.Background(), key, string(bytes), keyTTL).Err()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "insertRegionReport: cache operation failed - " + err.Error(),
		}
	}

	return nil
}

func (impl *CacheServiceImpl) GetRegionReport(region string, isQuarterReport bool, month, quarter, year int) (*models.RegionReportDB, *errors.Error) {

	if region == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "getRegionReport: cache missing details",
		}
	}

	var key string

	if isQuarterReport {
		key = fmt.Sprintf(constant.RedisOutletRegionReportKey, "", region, 0, quarter, year)
	} else {
		key = fmt.Sprintf(constant.RedisOutletRegionReportKey, "", region, month, 0, year)
	}

	data, err := impl.RedisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "getRegionReport: cache operation failed - " + err.Error(),
		}
	}

	if err == redis.Nil || data == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusNotFound,
			Code:       errors.UserNotFoundCode,
			Message:    errors.UserNotFound,
			Details:    "getRegionReport: cache no record found",
		}
	}

	report := &models.RegionReportDB{}

	err = objecthelper.UnmarshalJSON([]byte(data), &report)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "getRegionReport: cache unmarshal failed - " + err.Error(),
		}
	}

	return report, nil
}

func (impl *CacheServiceImpl) DelRegionReport(region string, isQuarterReport bool, month, quarter, year int) *errors.Error {

	if region == "" {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "delRegionReport: cache missing details",
		}
	}

	var key string

	if isQuarterReport {
		key = fmt.Sprintf(constant.RedisOutletRegionReportKey, "", region, 0, quarter, year)
	} else {
		key = fmt.Sprintf(constant.RedisOutletRegionReportKey, "", region, month, 0, year)
	}

	err := impl.RedisClient.Del(context.Background(), key).Err()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "delRegionReport: cache operation failed - " + err.Error(),
		}
	}

	return nil
}

func (impl *CacheServiceImpl) InsertUserSaleCountPerPlanReport(userId string, product string, isQuarterReport bool, month, quarter, year int, userReport interface{},
	keyTTL time.Duration) *errors.Error {

	if userReport == nil || userId == "" {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "InsertUserSaleCountPerPlanReport: cache missing details",
		}
	}

	bytes, err := objecthelper.MarshalJSON(userReport)
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "InsertUserSaleCountPerPlanReport: cache marshal failed - " + err.Error(),
		}
	}

	var key string

	if isQuarterReport {
		key = fmt.Sprintf(constant.RedisUserSaleCountPerPlanReportKey, userId, product, 0, quarter, year)
	} else {
		key = fmt.Sprintf(constant.RedisUserSaleCountPerPlanReportKey, userId, product, month, 0, year)
	}

	err = impl.RedisClient.Set(context.Background(), key, string(bytes), keyTTL).Err()
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "insertUserReport: cache operation failed - " + err.Error(),
		}
	}

	return nil
}

func (impl *CacheServiceImpl) GetUserSaleCountPerPlanReport(userId string, product string, isQuarterReport bool, month, quarter,
	year int) (interface{}, *errors.Error) {

	if userId == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "getUserSaleCountPerPlanReport: cache missing details",
		}
	}

	var key string

	if isQuarterReport {
		key = fmt.Sprintf(constant.RedisUserSaleCountPerPlanReportKey, userId, product, 0, quarter, year)
	} else {
		key = fmt.Sprintf(constant.RedisUserSaleCountPerPlanReportKey, userId, product, month, 0, year)
	}

	data, err := impl.RedisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "getUserSaleCountPerPlanReport: cache operation failed - " + err.Error(),
		}
	}

	if err == redis.Nil || data == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusNotFound,
			Code:       errors.UserNotFoundCode,
			Message:    errors.UserNotFound,
			Details:    "getUserSaleCountPerPlanReport: cache no record found",
		}
	}

	if product == constant.HBBCategory {
		report := &models.HbbUserSalesCountPerPlanDB{}
		err = objecthelper.UnmarshalJSON([]byte(data), &report)
		if err == nil {
			return report, nil
		}
	} else {
		report := &models.UserSalesCountPerPlanDB{}
		err = objecthelper.UnmarshalJSON([]byte(data), &report)
		if err == nil {
			return report, nil
		}
	}

	return nil, &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
		Details:    "getUserSaleCountPerPlanReport: cache unmarshal failed - " + err.Error(),
	}
}
