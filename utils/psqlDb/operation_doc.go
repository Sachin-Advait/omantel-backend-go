package psqldb

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (dbConn *DBImplementation) UpsertNotificationRecord(userId string, dbReq *models.NotificationDB) (string, *errors.Error) {

	err := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
		Details:    "invalid db arguments",
	}

	if dbReq == nil || dbReq.Id == "" {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var (
		dbErr  error
		result sql.Result
	)

	switch dbReq.IsUpdate {
	case false:
		query := `INSERT INTO user_notification(id, category, subject, body, file_type, file_link, created, modified, comment) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);`
		result, dbErr = dbConn.Operation.ExecContext(ctx, query, dbReq.Id, dbReq.Category, dbReq.Subject, dbReq.Body, dbReq.FileType, dbReq.FileLink, dbReq.CreatedAt, dbReq.ModifiedAt, "user - "+userId)
	default:
		query := `UPDATE user_notification SET category=$2, subject=$3, body=$4, file_type=$5, file_link=$6, modified=$7, comment=$8 WHERE id=$1 and deleted IS NULL`
		result, dbErr = dbConn.Operation.ExecContext(ctx, query, dbReq.Id, dbReq.Category, dbReq.Subject, dbReq.Body, dbReq.FileType, dbReq.FileLink, dbReq.ModifiedAt, "user - "+userId)
	}

	if dbErr != nil {
		if strings.Contains(dbErr.Error(), "duplicate") {
			return "", &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidDetailsCode,
				Message:    errors.InvalidDetails,
				Details:    fmt.Sprintf("query failed: %v", dbErr.Error()),
			}
		}

		err.Details = fmt.Sprintf("query failed: %v", dbErr.Error())
		return "", err
	}

	affectedRow, dbErr := result.RowsAffected()
	if dbErr != nil {
		err.Details = fmt.Sprintf("affected row failed: %v", dbErr.Error())
		return "", err
	}

	if affectedRow != 1 {
		err.Details = fmt.Sprintf("affected row invalid: %v", affectedRow)
		return "", err
	}

	return dbReq.Id, nil
}

func (dbConn *DBImplementation) FetchNotificationRecord(category string) ([]*models.NotificationDB, *errors.Error) {

	if category == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "invalid db arguments",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := "SELECT id, category, subject, body, file_type, file_link, created, modified FROM user_notification where deleted IS NULL and category=$1"

	rows, err := dbConn.Operation.QueryContext(ctx, query, category)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "FetchNotificationRecord query failed: " + err.Error(),
		}
	}

	response := make([]*models.NotificationDB, 0)

	for rows.Next() {

		var notif models.NotificationDB

		if err := rows.Scan(&notif.Id, &notif.Category, &notif.Subject, &notif.Body, &notif.FileType, &notif.FileLink, &notif.CreatedAt, &notif.ModifiedAt); err != nil {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "FetchNotificationRecord scan failed: " + err.Error(),
			}
		}

		response = append(response, &notif)
	}

	return response, nil
}
