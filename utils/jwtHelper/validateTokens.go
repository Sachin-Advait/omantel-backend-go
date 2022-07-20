package jwthelper

import (
	"errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func ValidateAccessToken(token, key string) (userID, userType, sourceID string, validate int, err error) {

	tokens, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {

		// if err == jwt.ErrSignatureInvalid {

		// }

		expiryErr, _ := err.(*jwt.ValidationError)

		if expiryErr.Errors == jwt.ValidationErrorExpired {
			return "", "", "", constant.TokenExpired, err
		}

		return "", "", "", constant.TokenInvalid, err
	}

	claims := tokens.Claims.(*Claims)

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) <= 0 {
		return "", "", "", constant.TokenExpired, errors.New("jwt token expired")
	}

	if claims.Token != constant.AccessTokenKey {
		return "", "", "", constant.TokenInvalid, errors.New("token invalid")
	}

	if claims.UserType == "" {
		return "", "", "", constant.TokenInvalid, errors.New("usertype invalid")
	}

	if claims.Subject == "" {
		return "", "", "", constant.TokenInvalid, errors.New("jwt token subject not found")
	}

	return claims.Subject, claims.UserType, claims.SourceID, constant.TokenValid, nil
}

func ValidateRefreshToken(token, key string) (userID, userType, sourceID string, validate int, err error) {

	tokens, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {

		// if err == jwt.ErrSignatureInvalid {

		// }

		expiryErr, _ := err.(*jwt.ValidationError)

		if expiryErr.Errors == jwt.ValidationErrorExpired {
			return "", "", "", constant.TokenExpired, err
		}

		return "", "", "", constant.TokenInvalid, err
	}

	claims := tokens.Claims.(*Claims)

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) <= 0 {
		return "", "", "", constant.TokenExpired, errors.New("jwt token expired")
	}

	if claims.Token != constant.RefreshTokenKey {
		return "", "", "", constant.TokenInvalid, errors.New("token invalid")
	}

	if claims.UserType == "" {
		return "", "", "", constant.TokenInvalid, errors.New("usertype invalid")
	}

	if claims.Subject == "" {
		return "", "", "", constant.TokenInvalid, errors.New("jwt token subject not found")
	}

	return claims.Subject, claims.UserType, claims.SourceID, constant.TokenValid, nil
}
