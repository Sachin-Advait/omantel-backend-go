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

func createNewRegionViaFile(ctx *gin.Context) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	if ctxHeader.Authorization == "" {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errors.UnauthorizedUser)).Stack().Msg("createNewRegionViaFile : unauthorized request")
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

	//  Ensure our file does not exceed 5MB
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 5*1024*1024)

	formFile, err := ctx.FormFile("region")
	if err != nil { // The file cannot be received.
		logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("createNewRegionViaFile : invalid file request")
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
		logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("createNewRegionViaFile : open file request failed")
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
		logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("createNewRegionViaFile : excelize reader failed")
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
		logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("createNewRegionViaFile : sheet1 not found")
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

	userRequests := make([]*models.NewRegionRequest, 0)

	for rows.Next() {

		if isFirstRow {
			firstRow, err := rows.Columns()
			if err != nil {
				logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("createNewRegionViaFile : first row error")
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
				case i == 1 && coulmn == "region_name":
					break
				case i == 2 && coulmn == "isUpdate":
					break
				default:
					logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("createNewRegionViaFile :  invalid column format")
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
			logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("createNewRegionViaFile : other rows error")
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

		userReq := &models.NewRegionRequest{}

		//LOOP:
		for i, colCell := range row {

			switch i {

			case 0:
				userReq.RegionCode = strings.ToUpper(strings.TrimSpace(colCell))
			case 1:
				userReq.RegionName = strings.ToUpper(strings.TrimSpace(colCell))
			case 2:
				userReq.IsUpdate, err = strconv.ParseBool(strings.TrimSpace(colCell))
			default:
				break
			}
		}

		userRequests = append(userRequests, userReq)
	}

	req, _ := objecthelper.MarshalJSON(userRequests)
	logger.FileLogger.CtxInfo(ctx).Interface(constant.BodyLogParam, fmt.Sprintf("%+v", string(req))).Msg("createNewRegionViaFile : req body")
	logger.LogChan <- &logger.Logging{
		UserID:     ctxHeader.XUserId,
		URL:        ctx.Request.URL.String(),
		Type:       "info",
		Msg:        "request body",
		Text:       fmt.Sprintf("%+v", string(req)),
		CreateTime: time.Now().Local(),
	}

	ctx.JSON(adapters.CreateNewRegion(ctx, userRequests))
}
