package impl

import (
	"context"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *CacheServiceImpl) InsertUseSessionToken(userId string, userSession *models.UserSessionToken, keyTTL time.Duration) *errors.Error {

	customError := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
	}

	if userId == "" || userSession == nil || userSession.AccessToken == "" || userSession.RefreshToken == "" {
		customError.Details = "insertUseSessionToken: cache missing details"
		return customError
	}

	valueBytes, err := objecthelper.MarshalJSON(userSession)
	if err != nil {
		customError.Details = "insertUseSessionToken: cache marshal failed - " + err.Error()
		return customError
	}

	err = impl.RedisClient.Set(context.Background(), userId, string(valueBytes), keyTTL).Err()
	if err != nil {
		customError.Details = "insertUseSessionToken: cache operation failed - " + err.Error()
		return customError
	}

	return nil
}

func (impl *CacheServiceImpl) GetUserSessionToken(userId string) (*models.UserSessionToken, *errors.Error) {

	customError := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
	}

	if userId == "" {
		customError.Details = "getUserSessionToken: cache missing details"
		return nil, customError
	}

	data, err := impl.RedisClient.Get(context.Background(), userId).Result()
	if err != nil && err != redis.Nil {
		customError.Details = "getUserSessionToken: cache operation failed - " + err.Error()
		return nil, customError
	}

	if err == redis.Nil || data == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusNotFound,
			Code:       errors.UserNotFoundCode,
			Message:    errors.UserNotFound,
			Details:    "getUserSessionToken: cache no record found",
		}
	}

	session := &models.UserSessionToken{}

	err = objecthelper.UnmarshalJSON([]byte(data), &session)
	if err != nil {
		customError.Details = "getUserSessionToken: cache unmarshal failed - " + err.Error()
		return nil, customError
	}

	return session, nil
}

func (impl *CacheServiceImpl) DelUserSessionToken(userId string) *errors.Error {

	customError := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
	}

	if userId == "" {
		customError.Details = "delUserSessionToken: cache missing details"
		return customError
	}

	err := impl.RedisClient.Del(context.Background(), userId).Err()
	if err != nil && err != redis.Nil {
		customError.Details = "delUserSessionToken: cache operation failed - " + err.Error()
		return customError
	}

	return nil
}

func (impl *CacheServiceImpl) IsSessionTokenValid(userId, sessionToken string) *errors.Error {

	customError := &errors.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       errors.InternalServerErrorCode,
		Message:    errors.InternalServerError,
	}

	if userId == "" || sessionToken == "" {
		customError.Details = "isSessionTokenValid: cache missing details"
		return customError
	}

	data, err := impl.RedisClient.Get(context.Background(), userId).Result()
	if err != nil && err != redis.Nil {
		customError.Details = "isSessionTokenValid: cache operation failed - " + err.Error()
		return customError
	}

	if err == redis.Nil || data == "" {
		return &errors.Error{
			StatusCode: http.StatusNotFound,
			Code:       errors.UserNotFoundCode,
			Message:    errors.UserNotFound,
			Details:    "isSessionTokenValid: cache no record found",
		}
	}

	session := &models.UserSessionToken{}

	err = objecthelper.UnmarshalJSON([]byte(data), &session)
	if err != nil {
		customError.Details = "isSessionTokenValid: cache unmarshal failed - " + err.Error()
		return customError
	}

	if session.AccessToken != sessionToken {
		customError.Details = "isSessionTokenValid: token invalid - " + sessionToken
		return customError
	}

	return nil
}
