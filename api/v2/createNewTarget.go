package v2

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gitlab.com/rendezvousPoint/incentive_calculator/adapters"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func createNewTargetViaFile(ctx *gin.Context) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	if ctxHeader.Authorization == "" {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errors.UnauthorizedUser)).Msg("createNewTargetViaFile : unauthorized request")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "unauthorized request",
			Text:       "auth token not found",
			CreateTime: time.Now().Local(),
		}

		ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
			Error:   errors.UnauthorizedUserCode,
			Message: errors.UnauthorizedUser,
		})
		return
	}

	if !ctxHeader.IsAdmin {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errors.UnauthorizedUser)).Msg("createNewTargetViaFile : not an admin")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "not admin",
			Code:       http.StatusForbidden,
			Text:       "not an admin",
			CreateTime: time.Now().Local(),
		}

		ctx.JSON(http.StatusForbidden, models.APIFailureResponse{
			Error:   errors.UnauthorizedUserCode,
			Message: errors.UnauthorizedUser,
		})
		return
	}

	//  Ensure our file does not exceed 5MB
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 5*1024*1024)

	formFile, err := ctx.FormFile("target")
	if err != nil { // The file cannot be received.
		logger.FileLogger.CtxError(ctx).Err(err).Msg("createNewTargetViaFile : invalid file request")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "invalid formFile",
			Text:       "invalid formFile",
			CreateTime: time.Now().Local(),
		}

		ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
			Error:   errors.InvalidReqParameterCode,
			Message: errors.InvalidReqParameters,
		})
		return
	}

	openedFile, err := formFile.Open()
	if err != nil { // The file cannot be received.
		logger.FileLogger.CtxError(ctx).Err(err).Msg("createNewTargetViaFile : open file request failed")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "open formFile failed",
			Text:       "open formFile failed",
			CreateTime: time.Now().Local(),
		}

		ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
			Error:   errors.InvalidReqParameterCode,
			Message: errors.InvalidReqParameters,
		})
		return
	}

	//_, _ = ioutil.ReadAll(openedFile)

	excelFile, err := excelize.OpenReader(openedFile)
	if err != nil { // The file cannot be received.
		logger.FileLogger.CtxError(ctx).Err(err).Msg("createNewTargetViaFile : excelize reader failed")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "invalid formFile format",
			Text:       "invalid formFile format",
			CreateTime: time.Now().Local(),
		}

		ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
			Error:   errors.InvalidReqParameterCode,
			Message: errors.InvalidReqParameters,
		})
		return
	}

	rows, err := excelFile.Rows("Sheet1")
	if err != nil {
		logger.FileLogger.CtxError(ctx).Err(err).Msg("createNewTargetViaFile : sheet1 not found")
		logger.LogChan <- &logger.Logging{
			UserID:     ctxHeader.XUserId,
			URL:        ctx.Request.URL.String(),
			Type:       "error",
			Msg:        "excel sheet missing",
			Text:       "excel sheet missing",
			CreateTime: time.Now().Local(),
		}

		ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
			Error:   errors.InvalidReqParameterCode,
			Message: errors.InvalidReqParameters,
		})
		return
	}

	isFirstRow := true

	userRequests := make([]*models.NewTargetRequest, 0)

	for rows.Next() {

		if isFirstRow {
			firstRow, err := rows.Columns()
			if err != nil {
				logger.FileLogger.CtxError(ctx).Err(err).Msg("createNewTargetViaFile : first row error")
				logger.LogChan <- &logger.Logging{
					UserID:     ctxHeader.XUserId,
					URL:        ctx.Request.URL.String(),
					Type:       "error",
					Msg:        "description row not found in sheet",
					Text:       "description row not found in sheet",
					CreateTime: time.Now().Local(),
				}

				ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
					Error:   errors.InvalidReqParameterCode,
					Message: errors.InvalidReqParameters,
				})
				return
			}

			for i, coulmn := range firstRow {

				switch {
				case i == 0 && coulmn == "region_code":
					break
				case i == 1 && coulmn == "outlet_code":
					break
				case i == 2 && coulmn == "year":
					break
				case i == 3 && coulmn == "month":
					break
				case i == 4 && coulmn == "baqati":
					break
				case i == 5 && coulmn == "hayyak":
					break
				case i == 6 && coulmn == "hbb":
					break
				case i == 7 && coulmn == "upgrade":
					break
				case i == 8 && coulmn == "isUpdate":
					break
				default:
					logger.FileLogger.CtxError(ctx).Err(err).Msg("createNewTargetViaFile :  invalid column format")
					logger.LogChan <- &logger.Logging{
						UserID:     ctxHeader.XUserId,
						URL:        ctx.Request.URL.String(),
						Type:       "error",
						Msg:        "invalid column format",
						Text:       "invalid column format",
						CreateTime: time.Now().Local(),
					}

					ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
						Error:   errors.InvalidReqParameterCode,
						Message: errors.InvalidReqParameters,
					})
					return
				}
			}

			isFirstRow = false
			continue
		}

		row, err := rows.Columns()
		if err != nil {
			logger.FileLogger.CtxError(ctx).Err(err).Msg("createNewTargetViaFile : other rows error")
			logger.LogChan <- &logger.Logging{
				UserID:     ctxHeader.XUserId,
				URL:        ctx.Request.URL.String(),
				Type:       "error",
				Msg:        "other rows not found in sheet",
				Text:       "other rows not found in sheet",
				CreateTime: time.Now().Local(),
			}

			ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
				Error:   errors.InvalidReqParameterCode,
				Message: errors.InvalidReqParameters,
			})
			return
		}

		userReq := &models.NewTargetRequest{}

		//LOOP:
		for i, colCell := range row {

			switch i {

			case 0:
				userReq.RegionCode = strings.ToUpper(strings.TrimSpace(colCell))
			case 1:
				userReq.OutletCode = strings.ToUpper(strings.TrimSpace(colCell))
			case 2:
				userReq.Year, err = strconv.Atoi(strings.TrimSpace(colCell))
			case 3:
				userReq.Month, err = strconv.Atoi(strings.TrimSpace(colCell))
			case 4:
				userReq.Baqati, err = strconv.ParseFloat(strings.TrimSpace(colCell), 64)
			case 5:
				userReq.Hayyak, err = strconv.ParseFloat(strings.TrimSpace(colCell), 64)
			case 6:
				userReq.HBB, err = strconv.ParseFloat(strings.TrimSpace(colCell), 64)
			case 7:
				userReq.Upgrade, err = strconv.ParseFloat(strings.TrimSpace(colCell), 64)
			case 8:
				userReq.IsUpdate, err = strconv.ParseBool(strings.TrimSpace(colCell))
			default:
				break
			}
		}

		userRequests = append(userRequests, userReq)
	}

	req, _ := objecthelper.MarshalJSON(userRequests)
	logger.FileLogger.CtxInfo(ctx).Interface(logger.BodyLogParam, fmt.Sprintf("%+v", string(req))).Msg("createNewTargetViaFile : req body")
	logger.LogChan <- &logger.Logging{
		UserID:     ctxHeader.XUserId,
		URL:        ctx.Request.URL.String(),
		Type:       "info",
		Msg:        "request body",
		Text:       fmt.Sprintf("%+v", string(req)),
		CreateTime: time.Now().Local(),
	}

	ctx.JSON(adapters.CreateNewTarget(ctx, userRequests))
}
