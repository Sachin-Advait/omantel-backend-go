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
	dblogger "gitlab.com/rendezvousPoint/incentive_calculator/utils/logger/dbLogger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func UpsertVOCUser(ctx *gin.Context) (int, interface{}) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	userReq, err := upsertVOCUserLoadFromJsonBody(ctx)
	if err != nil {
		userReq, err = upsertVOCUserLoadFromFileBody(ctx)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Err(err).Msg("upsertVOCUser : invalid request param")
			dblogger.PrintDBLogs(&dblogger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "request binding", Code: http.StatusBadRequest,
				Text: err.Error(), CreateTime: time.Now().Local()})

			return http.StatusBadRequest, models.APIFailureResponse{
				Status:  false,
				Error:   errors.InvalidReqParameters,
				Message: errors.RequestDecodeError,
			}
		}
	}

	req, _ := objecthelper.MarshalJSON(userReq)
	logger.FileLogger.CtxInfo(ctx).Interface(logger.BodyLogParam, fmt.Sprintf("%+v", string(req))).Msg("upsertVOCUser : req body")
	dblogger.PrintDBLogs(&dblogger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "trace", Msg: "request body",
		Text: fmt.Sprintf("%+v", string(req)), CreateTime: time.Now().Local()})

	if err := models.ValidateUpsertVOCUser(userReq); err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("upsertVOCUser : req validation failed")
		dblogger.PrintDBLogs(&dblogger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "invalid req params",
			Text: err.Error(), CreateTime: time.Now().Local()})

		return http.StatusBadRequest, &models.APIFailureResponse{
			Status:  false,
			Error:   errors.InvalidReqParameters,
			Message: err.Error(),
		}
	}

	failedResp, customErr := adminAPIProvider.UpsertVOCUser(ctx, userReq)
	if customErr != nil {
		logger.FileLogger.CtxError(ctx).Err(errors.New(customErr.Details.(string))).Msg("upsertVOCUser : provider failed")
		dblogger.PrintDBLogs(&dblogger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "error", Msg: "provider failed",
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

	logger.FileLogger.CtxInfo(ctx).Str(logger.ResponseParam, errors.SuccessInfoMessage).Msg("upsertVOCUser : success")
	dblogger.PrintDBLogs(&dblogger.Logging{UserID: ctxHeader.XUserId, URL: ctx.Request.URL.String(), Type: "info", Msg: "success",
		Text: errors.SuccessInfoMessage, CreateTime: time.Now().Local()})

	return http.StatusOK, models.APISuccessResponse{
		Status: true,
		Data:   errors.SuccessInfoMessage,
	}
}

func upsertVOCUserLoadFromJsonBody(ctx *gin.Context) ([]*models.CreateVOCUserRequest, error) {

	userReq := &models.CreateVOCUserRequest{}
	if err := ctx.ShouldBindJSON(userReq); err != nil {
		return nil, err
	}

	return []*models.CreateVOCUserRequest{userReq}, nil
}

func upsertVOCUserLoadFromFileBody(ctx *gin.Context) ([]*models.CreateVOCUserRequest, error) {

	//  Ensure our file does not exceed 5MB
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 5*1024*1024)

	formFile, err := ctx.FormFile("voc_user")
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
		userRequests []*models.CreateVOCUserRequest
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
				case i == 0 && strings.ToUpper(strings.TrimSpace(coulmn)) == "AGENT_ID":
					continue
				case i == 1 && strings.ToUpper(strings.TrimSpace(coulmn)) == "PERCENTAGE":
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

		userReq := &models.CreateVOCUserRequest{}

		//LOOP:
		for i, colCell := range row {

			switch i {

			case 0:
				userReq.UserId = strings.TrimSpace(colCell)
				if userReq.UserId == "" {
					return nil, fmt.Errorf("userId - row:%v, err:%v", index, err.Error())
				}
			case 1:
				userReq.AchievedPercent, err = strconv.ParseFloat(strings.TrimSpace(colCell), 64)
				if err != nil {
					return nil, fmt.Errorf("achievedPercent - row:%v, err:%v", index, err.Error())
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
