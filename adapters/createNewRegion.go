package adapters

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func CreateNewRegion(ctx *gin.Context) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	userReq, err := regionPayloadFromJsonBody(ctx)
	if err != nil {
		userReq, err = regionPayloadFromFileBody(ctx)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Err(err).Msg("upsertNewRegion : invalid request param")
			logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "request binding", Code: http.StatusBadRequest,
				Text: err.Error(), CreateTime: time.Now().Local()})

			return http.StatusBadRequest, models.APIFailureResponse{
				Error:   errors.RequestDecodeErrorCode,
				Message: errors.RequestDecodeError,
			}
		}
	}

	req, _ := objecthelper.MarshalJSON(userReq)
	logger.FileLogger.CtxInfo(ctx).Interface(logger.BodyLogParam, fmt.Sprintf("%+v", string(req))).Msg("upsertNewRegion : req body")
	logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "trace", Msg: "request body",
		Text: fmt.Sprintf("%+v", string(req)), CreateTime: time.Now().Local()})

	if err := models.ValidateUpsertRegion(userReq); err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("upsertNewRegion : req validation failed")
		logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "invalid req params",
			Text: err.Error(), CreateTime: time.Now().Local()})

		return http.StatusBadRequest, &models.APIFailureResponse{
			Error:   errors.InvalidReqParameters,
			Message: err.Error(),
		}
	}

	failedResp, customErr := theProvider.CreateNewRegion(ctx, userReq)
	if customErr != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(customErr.Details.(string))).Msg("upsertNewRegion : provider failed")
		logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "provider failed",
			Text: customErr.Details.(string), CreateTime: time.Now().Local()})

		if len(failedResp) != 0 {
			return http.StatusOK, failedResp
		}

		return customErr.StatusCode, &models.APIFailureResponse{
			Error: customErr.Message,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, errors.SuccessInfoMessage).Msg("upsertNewRegion : success")
	logger.PrintDBLogs(&logger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "info", Msg: "success",
		Text: errors.SuccessInfoMessage, CreateTime: time.Now().Local()})

	return http.StatusOK, models.APISuccessResponse{
		Status: true,
		Data:   errors.SuccessInfoMessage,
	}
}

func regionPayloadFromJsonBody(ctx *gin.Context) ([]*models.CreateRegionRequest, error) {

	userReq := &models.CreateRegionRequest{}
	if err := ctx.ShouldBindJSON(userReq); err != nil {
		return nil, err
	}

	return []*models.CreateRegionRequest{userReq}, nil
}

func regionPayloadFromFileBody(ctx *gin.Context) ([]*models.CreateRegionRequest, error) {

	//  Ensure our file does not exceed 5MB
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 5*1024*1024)

	formFile, err := ctx.FormFile("region")
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
		userRequests []*models.CreateRegionRequest
	)

OUTERLOOP:
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
				case i == 2 && coulmn == "quarter":
					continue
				case i == 3 && coulmn == "year":
					continue
				case i == 4 && coulmn == "isUpdate":
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

		userReq := &models.CreateRegionRequest{}

		//LOOP:
		for i, colCell := range row {

			switch i {

			case 0:
				userReq.RegionCode = strings.ToUpper(strings.TrimSpace(colCell))
				if userReq.RegionCode == "" {
					continue OUTERLOOP
				}
			case 1:
				userReq.RegionName = strings.ToUpper(strings.TrimSpace(colCell))
			case 2:
				userReq.Quarter, err = strconv.Atoi(strings.TrimSpace(colCell))
				if err != nil {
					return nil, fmt.Errorf("quarter - row:%v, err:%v", index, err.Error())
				}
			case 3:
				userReq.Year, err = strconv.Atoi(strings.TrimSpace(colCell))
				if err != nil {
					return nil, fmt.Errorf("year - row:%v, err:%v", index, err.Error())
				}
			case 4:
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
