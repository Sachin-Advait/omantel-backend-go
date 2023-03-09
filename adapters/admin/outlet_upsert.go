package admin

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func UpsertOutlet(ctx *gin.Context) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	userReq, err := outletPayloadFromJsonBody(ctx)
	if err != nil {
		userReq, err = outletPayloadFromFileBody(ctx)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Err(err).Msg("upsertNewOutlet : invalid request param")
			logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "request binding", Code: http.StatusBadRequest,
				Text: err.Error(), CreateTime: time.Now().Local()})

			return http.StatusBadRequest, models.APIFailureResponse{
				Status:  false,
				Error:   errors.InvalidReqParameters,
				Message: errors.RequestDecodeError,
			}
		}
	}

	req, _ := objecthelper.MarshalJSON(userReq)
	logger.FileLogger.CtxInfo(ctx).Interface(logger.BodyLogParam, fmt.Sprintf("%+v", string(req))).Msg("upsertNewOutlet : req body")
	logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "trace", Msg: "request body",
		Text: fmt.Sprintf("%+v", string(req)), CreateTime: time.Now().Local()})

	if err := models.ValidateUpsertOutlet(userReq); err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("upsertNewOutlet : req validation failed")
		logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "invalid req params",
			Text: err.Error(), CreateTime: time.Now().Local()})

		return http.StatusBadRequest, &models.APIFailureResponse{
			Status:  false,
			Error:   errors.InvalidReqParameters,
			Message: err.Error(),
		}
	}

	failedResp, customErr := adminAPIProvider.UpsertOutlet(ctx, userReq)
	if customErr != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(customErr.Details.(string))).Msg("upsertNewOutlet : provider failed")
		logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "provider failed",
			Text: customErr.Details.(string), CreateTime: time.Now().Local()})

		if len(failedResp) != 0 {
			return http.StatusOK, models.APIResponse{
				Status:  false,
				Message: "PARTIAL SUCCESS AND FAILURE",
				Data:    failedResp,
			}
		}

		return customErr.StatusCode, &models.APIFailureResponse{
			Status:  false,
			Error:   customErr.Err,
			Message: customErr.Message,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, errors.SuccessInfoMessage).Msg("upsertNewOutlet : success")
	logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "info", Msg: "success",
		Text: errors.SuccessInfoMessage, CreateTime: time.Now().Local()})

	return http.StatusOK, models.APISuccessResponse{
		Status: true,
		Data:   errors.SuccessInfoMessage,
	}
}

func outletPayloadFromJsonBody(ctx *gin.Context) ([]*models.CreateOutletRequest, error) {

	userReq := &models.CreateOutletRequest{}
	if err := ctx.ShouldBindJSON(userReq); err != nil {
		return nil, err
	}

	return []*models.CreateOutletRequest{userReq}, nil
}

func outletPayloadFromFileBody(ctx *gin.Context) ([]*models.CreateOutletRequest, error) {

	//  Ensure our file does not exceed 5MB
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 5*1024*1024)

	formFile, err := ctx.FormFile("outlet")
	if err != nil {
		return nil, fmt.Errorf("formFile - %v", err.Error())
	}

	openedFile, err := formFile.Open()
	if err != nil {
		return nil, fmt.Errorf("openFile - %v", err.Error())
	}

	excelFile, err := excelize.OpenReader(openedFile)
	if err != nil {
		return nil, fmt.Errorf("openReader - %v", err.Error())
	}

	rows, err := excelFile.Rows("Sheet 1")
	if err != nil {
		return nil, fmt.Errorf("rows - %v", err.Error())
	}

	var (
		index        int
		isFirstRow   = true
		userRequests []*models.CreateOutletRequest
	)

	for rows.Next() {

		index++

		if isFirstRow {
			firstRow, err := rows.Columns()
			if err != nil {
				return nil, fmt.Errorf("firstColumn - %v", err.Error())
			}

			//firstRow = firstRow[:9]

			for i, coulmn := range firstRow {

				switch {
				case i == 0 && coulmn == "region_code":
					continue
				case i == 1 && coulmn == "region_name":
					continue
				case i == 2 && coulmn == "outlet_code":
					continue
				case i == 3 && coulmn == "outlet_name":
					continue
				case i == 4 && coulmn == "lat":
					continue
				case i == 5 && coulmn == "lang":
					continue
				case i == 6 && coulmn == "quarter":
					continue
				case i == 7 && coulmn == "year":
					continue
				case i == 8 && coulmn == "isUpdate":
					continue
				default:
					return nil, errors.New("invalid column format")
				}
			}

			isFirstRow = false
			continue
		}

		row, err := rows.Columns()
		if err != nil {
			return nil, fmt.Errorf("columns - %v", err.Error())
		}

		userReq := &models.CreateOutletRequest{}

		//LOOP:
		for i, colCell := range row {

			switch i {

			case 0:
				userReq.RegionCode = strings.ToUpper(strings.TrimSpace(colCell))
			case 1:
				userReq.RegionName = strings.ToUpper(strings.TrimSpace(colCell))
			case 2:
				userReq.OutletCode = strings.ToUpper(strings.TrimSpace(colCell))
			case 3:
				userReq.OutletName = strings.ToUpper(strings.TrimSpace(colCell))
			case 4:
				userReq.Lat, err = strconv.ParseFloat(strings.TrimSpace(colCell), 64)
				if err != nil {
					return nil, fmt.Errorf("latitude - row:%v, err:%v", index, err.Error())
				}
			case 5:
				userReq.Lang, err = strconv.ParseFloat(strings.TrimSpace(colCell), 64)
				if err != nil {
					return nil, fmt.Errorf("longitude - row:%v, err:%v", index, err.Error())
				}
			case 6:
				userReq.Quarter, err = strconv.Atoi(strings.TrimSpace(colCell))
				if err != nil {
					return nil, fmt.Errorf("quarter - row:%v, err:%v", index, err.Error())
				}
			case 7:
				userReq.Year, err = strconv.Atoi(strings.TrimSpace(colCell))
				if err != nil {
					return nil, fmt.Errorf("year - row:%v, err:%v", index, err.Error())
				}
			case 8:
				userReq.IsUpdate, err = strconv.ParseBool(strings.TrimSpace(colCell))
				if err != nil {
					return nil, fmt.Errorf("isUpdate - row:%v, err:%v", index, err.Error())
				}
			}
		}

		userRequests = append(userRequests, userReq)
	}

	if len(userRequests) == 0 {
		return nil, errors.New("no record")
	}

	return userRequests, nil
}
