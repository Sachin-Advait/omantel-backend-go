package impl

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (dbConn *DbServiceImpl) InsertLoginDetails(login *models.LoginUsers) *errors.Error {

	// logger.FileLogger.Info().Interface("clientID", staffId).Time("start", time.Now()).Msg("InsertLoginDetails start")
	// defer logger.FileLogger.Info().Interface("clientID", staffId).Time("end", time.Now()).Msg("InsertLoginDetails done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	_, err := dbConn.DBClient.QueryContext(ctx, "INSERT INTO user_login_table(emp_id, region, login) VALUES ($1,$2,$3);",
		&login.StaffId,
		&login.StaffRegion,
		&login.LoginTime,
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

	var (
		err   error
		count int
	)

	if region != "" {
		err = dbConn.DBClient.QueryRowContext(ctx, "SELECT COUNT(*) FROM (SELECT DISTINCT emp_id FROM user_login_table where "+
			"region=$1 AND login >= $2 AND login <= $3) AS total;", region, start, end).Scan(&count)
	} else {
		err = dbConn.DBClient.QueryRowContext(ctx, "SELECT COUNT(*) FROM (SELECT DISTINCT emp_id FROM user_login_table where "+
			"login >= $1 AND login <= $2) AS total;", start, end).Scan(&count)
	}

	if err != nil {
		return count, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchLoginCount: query failed-" + err.Error(),
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
		rows, dbErr = dbConn.DBClient.QueryContext(ctx, "SELECT emp_id, region, login FROM user_login_table where region=$1 AND login >= $2 AND login <= $3 "+
			"order by login desc LIMIT $4 OFFSET $5 ;", region, start, end, size, offset)
	} else {
		rows, dbErr = dbConn.DBClient.QueryContext(ctx, "SELECT emp_id, region, login FROM user_login_table where login >= $1 AND login <= $2 "+
			"order by login desc LIMIT $3 OFFSET $4 ;", start, end, size, offset)
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
		if dbErr := rows.Scan(&user.StaffId, &user.StaffRegion, &user.LoginTime); dbErr != nil {
			err.Details = "FetchLoginDetails: scan failed-" + dbErr.Error()
			return nil, err
		}
		response = append(response, &user)
	}

	return response, nil
}
