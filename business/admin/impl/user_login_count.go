package impl

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *AdminImpl) GetUserLoginCount(ctx *gin.Context, region string, start time.Time, end time.Time) (*models.UserLoginCountResponse, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	if err := impl.CacheService.IsSessionTokenValid(ctxHeader.XUserId, ctxHeader.Authorization); err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusUnauthorized,
			Code:       errors.UnauthorizedUserCode,
			Message:    errors.UnauthorizedUser,
			Details:    err.Details,
		}
	}

	quarter := getQuarterFromMonth(int(start.Month()))

	var (
		err        *errors.Error
		userCount  int
		loginCount int
	)

	var wg sync.WaitGroup
	wg.Add(2)

	go func(w *sync.WaitGroup) {
		defer w.Done()
		count, userErr := impl.DbService.FetchDistinctUserCountFromOutletAndRegion(region, "", quarter, start.Year())
		if userErr != nil {
			err = &errors.Error{StatusCode: userErr.StatusCode, Code: userErr.Code, Err: userErr.Err, Message: userErr.Message, Details: userErr.Details}
			return
		}
		userCount = count
	}(&wg)

	go func(w *sync.WaitGroup) {
		defer w.Done()
		count, countErr := impl.DbService.FetchLoginCount(region, start, end)
		if countErr != nil {
			err = &errors.Error{StatusCode: countErr.StatusCode, Code: countErr.Code, Err: countErr.Err, Message: countErr.Message, Details: countErr.Details}
			return
		}
		loginCount = count
	}(&wg)

	wg.Done()

	if err != nil {
		return nil, err
	}

	response := &models.UserLoginCountResponse{
		TotalUsers:    userCount,
		LoginCount:    loginCount,
		NonLoginCount: userCount - loginCount,
	}

	response.LoginPercentage = (float64(loginCount) * 100) / float64(userCount)
	response.NonLoginPercentage = 100 - response.LoginPercentage
	return response, nil
}
