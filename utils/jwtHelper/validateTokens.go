package jwthelper

import (
	"errors"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/constant"

	"github.com/dgrijalva/jwt-go"
)

func ValidateAccessToken(token, key string) (userID string, validate int, err error) {

	tokens, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {

		// if err == jwt.ErrSignatureInvalid {

		// }

		expiryErr, _ := err.(*jwt.ValidationError)

		if expiryErr.Errors == jwt.ValidationErrorExpired {
			return "", constant.TokenExpired, err
		}

		return "", constant.TokenInvalid, err
	}

	claims := tokens.Claims.(*Claims)

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) <= 0 {
		return "", constant.TokenExpired, errors.New("jwt token expired")
	}

	if claims.Token != constant.AccessTokenKey {
		return "", constant.TokenInvalid, errors.New("token invalid")
	}

	if claims.Subject == "" {
		return "", constant.TokenInvalid, errors.New("jwt token subject not found")
	}

	return claims.Subject, constant.TokenValid, nil
}

func ValidateRefreshToken(token, key string) (userID string, validate int, err error) {

	tokens, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {

		// if err == jwt.ErrSignatureInvalid {

		// }

		expiryErr, _ := err.(*jwt.ValidationError)

		if expiryErr.Errors == jwt.ValidationErrorExpired {
			return "", constant.TokenExpired, err
		}

		return "", constant.TokenInvalid, err
	}

	claims := tokens.Claims.(*Claims)

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) <= 0 {
		return "", constant.TokenExpired, errors.New("jwt token expired")
	}

	if claims.Token != constant.RefreshTokenKey {
		return "", constant.TokenInvalid, errors.New("token invalid")
	}

	if claims.Subject == "" {
		return "", constant.TokenInvalid, errors.New("jwt token subject not found")
	}

	return claims.Subject, constant.TokenValid, nil
}
