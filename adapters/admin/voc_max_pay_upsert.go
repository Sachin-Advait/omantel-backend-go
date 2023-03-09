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

func UpsertVOCMaxPay(ctx *gin.Context) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	userReq, err := upsertVOCMaxPayLoadFromJsonBody(ctx)
	if err != nil {
		userReq, err = upsertVOCMaxPayLoadFromFileBody(ctx)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Err(err).Msg("upsertVOCMaxPay : invalid request param")
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
	logger.FileLogger.CtxInfo(ctx).Interface(logger.BodyLogParam, fmt.Sprintf("%+v", string(req))).Msg("upsertVOCMaxPay : req body")
	logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "trace", Msg: "request body",
		Text: fmt.Sprintf("%+v", string(req)), CreateTime: time.Now().Local()})

	if err := models.ValidateUpsertVOCMaxPay(userReq); err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("upsertVOCMaxPay : req validation failed")
		logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "invalid req params",
			Text: err.Error(), CreateTime: time.Now().Local()})

		return http.StatusBadRequest, &models.APIFailureResponse{
			Status:  false,
			Error:   errors.InvalidReqParameters,
			Message: err.Error(),
		}
	}

	failedResp, customErr := adminAPIProvider.UpsertVOCMaxPay(ctx, userReq)
	if customErr != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(customErr.Details.(string))).Msg("upsertVOCMaxPay : provider failed")
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

	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, errors.SuccessInfoMessage).Msg("upsertVOCMaxPay : success")
	logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "info", Msg: "success",
		Text: errors.SuccessInfoMessage, CreateTime: time.Now().Local()})

	return http.StatusOK, models.APISuccessResponse{
		Status: true,
		Data:   errors.SuccessInfoMessage,
	}
}

func upsertVOCMaxPayLoadFromJsonBody(ctx *gin.Context) ([]*models.CreateVOCMaxPayRequest, error) {

	userReq := &models.CreateVOCMaxPayRequest{}
	if err := ctx.ShouldBindJSON(userReq); err != nil {
		return nil, err
	}

	return []*models.CreateVOCMaxPayRequest{userReq}, nil
}

func upsertVOCMaxPayLoadFromFileBody(ctx *gin.Context) ([]*models.CreateVOCMaxPayRequest, error) {

	//  Ensure our file does not exceed 5MB
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 5*1024*1024)

	formFile, err := ctx.FormFile("voc_max_pay")
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
		userRequests []*models.CreateVOCMaxPayRequest
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
				case i == 0 && strings.ToUpper(strings.TrimSpace(coulmn)) == "POSITION":
					continue
				case i == 1 && strings.ToUpper(strings.TrimSpace(coulmn)) == "MAX_VOC_PAY":
					continue
				case i == 2 && strings.ToUpper(strings.TrimSpace(coulmn)) == "YEAR":
					continue
				case i == 3 && strings.ToUpper(strings.TrimSpace(coulmn)) == "ISUPDATE":
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

		userReq := &models.CreateVOCMaxPayRequest{}

		//LOOP:
		for i, colCell := range row {

			switch i {

			case 0:
				userReq.Position = strings.TrimSpace(colCell)
				_, ok := constant.EmployeePositionMap[userReq.Position]
				if !ok {
					return nil, fmt.Errorf("position - row:%v, err:%v", index, err.Error())
				}
			case 1:
				userReq.MaxPay, err = strconv.ParseFloat(strings.TrimSpace(colCell), 64)
				if err != nil {
					return nil, fmt.Errorf("maxPay - row:%v, err:%v", index, err.Error())
				}
			case 2:
				userReq.Year, err = strconv.Atoi(strings.TrimSpace(colCell))
				if err != nil {
					return nil, fmt.Errorf("year - row:%v, err:%v", index, err.Error())
				}
			case 3:
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
