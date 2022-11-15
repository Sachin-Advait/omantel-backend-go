package adapters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func CreateNewProduct(ctx *gin.Context, req []*models.NewProductRequest) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	for _, eachRecord := range req {

		if eachRecord.ProductCode == "" || eachRecord.ProductName == "" || eachRecord.Weightage == 0 {
			logger.FileLogger.CtxError(ctx).Err(errors.New("invalid product param")).Stack().Msg("CreateNewProduct : invalid req params")
			logger.LogChan <- &logger.Logging{
				UserID:     ctxHeader.XUserId,
				URL:        ctx.Request.URL.String(),
				Type:       "error",
				Msg:        "invalid req params",
				Text:       "invalid req param",
				CreateTime: time.Now().Local(),
			}

			return http.StatusBadRequest, &models.APIFailureResponse{
				Error:   errors.InvalidReqParameters,
				Message: "Invalid product params",
			}
		}
	}

	for i := 0; i < len(req); i++ {
		req[i].ProductCode = strings.ToUpper(strings.TrimSpace(req[i].ProductCode))
		req[i].ProductName = strings.ToUpper(strings.TrimSpace(req[i].ProductName))
	}

	failedResp, errs := theProvider.CreateNewProduct(ctx, req)

	switch ctxHeader.XAPIVersion {
	case "2":
		if failedResp != nil {
			resp, _ := json.Marshal(failedResp)
			logger.FileLogger.CtxError(ctx).Str(constant.ResponseParam, string(resp)).Stack().Msg("CreateNewProduct V2: provider failed")
			logger.LogChan <- &logger.Logging{
				UserID:     ctxHeader.XUserId,
				URL:        ctx.Request.URL.String(),
				Type:       "error",
				Msg:        "provider failed",
				Text:       string(resp),
				CreateTime: time.Now().Local(),
			}

			return http.StatusOK, failedResp
		}
	default:
		if errs != nil {
			logger.FileLogger.CtxError(ctx).Err(errors.New(errs.Details.(string))).Stack().Msg("CreateNewProduct : provider failed")
			logger.LogChan <- &logger.Logging{
				UserID:     ctxHeader.XUserId,
				URL:        ctx.Request.URL.String(),
				Type:       "error",
				Msg:        "provider failed",
				Text:       errs.Details.(string),
				CreateTime: time.Now().Local(),
			}

			return errs.StatusCode, &models.APIFailureResponse{
				Error: errs.Message,
			}
		}
	}

	resp := models.APISuccessResponse{
		Status: true,
		Data:   errors.SuccessInfoMessage,
	}

	logger.FileLogger.CtxInfo(ctx).Str(constant.ResponseParam, fmt.Sprintf("%+v", resp)).Msg("CreateNewProduct : success")
	logger.LogChan <- &logger.Logging{
		UserID:     ctxHeader.XUserId,
		URL:        ctx.Request.URL.String(),
		Type:       "info",
		Msg:        "success",
		Text:       fmt.Sprintf("%+v", resp),
		CreateTime: time.Now().Local(),
	}

	return http.StatusOK, resp
}
