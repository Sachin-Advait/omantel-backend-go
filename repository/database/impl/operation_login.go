package impl

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (dbConn *DbServiceImpl) InsertLoginDetails(login *models.UserLoginDetailsDB) *errors.Error {

	// logger.FileLogger.Info().Interface("clientID", staffId).Time("start", time.Now()).Msg("InsertLoginDetails start")
	// defer logger.FileLogger.Info().Interface("clientID", staffId).Time("end", time.Now()).Msg("InsertLoginDetails done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	_, err := dbConn.DBClient.QueryContext(ctx, "INSERT INTO user_login_table(emp_id, month, quarter, year, login) VALUES ($1,$2,$3,$4,$5);",
		login.UserId,
		login.Month,
		login.Quarter,
		login.Year,
		login.LoginTime,
	)

	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "InsertLoginDetails: insert failed - " + err.Error(),
		}
	}

	return nil
}

func (dbConn *DbServiceImpl) FetchLoginDetails(userID string, isQuarterReport bool, quarter, month, year int) ([]time.Time, *errors.Error) {

	// logger.FileLogger.Info().Str("clientID", userID).Time("start", time.Now()).Msg("FetchLoginDetails start")
	// defer logger.FileLogger.Info().Str("clientID", userID).Time("end", time.Now()).Msg("FetchLoginDetails done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var (
		rows  *sql.Rows
		dbErr error
		err   = &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
		}
	)

	switch {
	case isQuarterReport:
		rows, dbErr = dbConn.DBClient.QueryContext(ctx, "select login from user_login_table where emp_id=$1 and quarter=$2 and year=$3", userID, quarter, year)
	default:
		rows, dbErr = dbConn.DBClient.QueryContext(ctx, "select login from user_login_table where emp_id=$1 and month=$2 and year=$3", userID, month, year)
	}

	if dbErr != nil {
		if dbErr == sql.ErrNoRows {
			return nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidDetailsCode,
				Err:        errors.InvalidDetails,
				Message:    "NO RECORDS FOUND",
				Details:    "FetchLoginDetails: query failed-" + dbErr.Error(),
			}
		}

		err.Details = "FetchLoginDetails: query failed-" + dbErr.Error()
		return nil, err
	}

	response := make([]time.Time, 0)

	for rows.Next() {
		var loginTime time.Time
		if dbErr := rows.Scan(&loginTime); dbErr != nil {
			err.Details = "FetchLoginDetails: scan failed-" + dbErr.Error()
			return nil, err
		}
		response = append(response, loginTime)
	}

	return response, nil
}
