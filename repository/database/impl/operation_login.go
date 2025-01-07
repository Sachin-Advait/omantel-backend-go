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

func (dbConn *DbServiceImpl) FetchLoginCount(region string, start, end time.Time) (int, *errors.Error) {

	// logger.FileLogger.Info().Str("clientID", userID).Time("start", time.Now()).Msg("FetchLoginCount start")
	// defer logger.FileLogger.Info().Str("clientID", userID).Time("end", time.Now()).Msg("FetchLoginCount done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	startDate := start.Format("2006-01-02")
	endDate := start.Format("2006-01-02")

	var (
		err   error
		count int
	)

	if region != "" {
		err = dbConn.DBClient.QueryRowContext(ctx, "SELECT COUNT(*) FROM (SELECT DISTINCT emp_id FROM user_login_table where "+
			"region=$1 AND login >= '$2 00:00:00'::timestamp AND login < '$3 23:59:59'::timestamp) AS total;", region, startDate, endDate).Scan(&count)
	} else {
		err = dbConn.DBClient.QueryRowContext(ctx, "SELECT COUNT(*) FROM (SELECT DISTINCT emp_id FROM user_login_table where "+
			"login >= '$1 00:00:00'::timestamp AND login < '$2 23:59:59'::timestamp) AS total;", startDate, endDate).Scan(&count)
	}

	if err != nil {
		return count, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
		}
	}

	return count, nil
}

func (dbConn *DbServiceImpl) FetchLoginDetails(page int, size int, region string, start, end time.Time) ([]*models.LoginUsers, *errors.Error) {

	// logger.FileLogger.Info().Str("clientID", userID).Time("start", time.Now()).Msg("FetchLoginDetails start")
	// defer logger.FileLogger.Info().Str("clientID", userID).Time("end", time.Now()).Msg("FetchLoginDetails done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	offset := (page - 1) * size
	if offset < 0 {
		offset = 0
	}
	startDate := start.Format("2006-01-02")
	endDate := start.Format("2006-01-02")

	var (
		rows  *sql.Rows
		dbErr error
		err   = &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
		}
	)

	if region != "" {
		rows, dbErr = dbConn.DBClient.QueryContext(ctx, "SELECT emp_id, login FROM user_login_table where region=$1 AND login >= '$2 00:00:00'::timestamp AND "+
			"login < '$3 23:59:59'::timestamp) order by login desc LIMIT $4 OFFSET $5 ;", region, startDate, endDate, size, offset)
	} else {
		rows, dbErr = dbConn.DBClient.QueryContext(ctx, "SELECT emp_id, login FROM user_login_table where login >= '$1 00:00:00'::timestamp AND "+
			"login < '$2 23:59:59'::timestamp) order by login desc LIMIT $4 OFFSET $5 ;", startDate, endDate, size, offset)
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

	response := make([]*models.LoginUsers, 0)

	for rows.Next() {
		var user models.LoginUsers
		if dbErr := rows.Scan(&user.StaffId, &user.LoginTime); dbErr != nil {
			err.Details = "FetchLoginDetails: scan failed-" + dbErr.Error()
			return nil, err
		}
		response = append(response, &user)
	}

	return response, nil
}
