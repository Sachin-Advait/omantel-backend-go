package jwthelper

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
)

func CreateAccessToken(id string, expiry time.Duration, key string) (string, time.Time, error) {

	issuedAt := time.Now()
	exp := issuedAt.Add(expiry)

	atClaims := jwt.MapClaims{}
	atClaims["iat"] = issuedAt.Unix()
	atClaims["sub"] = id
	atClaims["token"] = constant.AccessTokenKey
	atClaims["exp"] = exp.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)

	accessToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken, exp, nil
}

func CreateRefreshToken(id string, expiry time.Duration, key string) (string, time.Time, error) {

	issuedAt := time.Now()
	exp := issuedAt.Add(expiry)

	atClaims := jwt.MapClaims{}
	atClaims["iat"] = issuedAt.Unix()
	atClaims["sub"] = id
	atClaims["token"] = constant.RefreshTokenKey
	atClaims["exp"] = exp.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)

	refreshToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", time.Time{}, err
	}

	return refreshToken, exp, nil
}

func RefreshAccessToken(refreshToken string, expiry time.Duration, key string) (string, time.Time, error) {

	tokens, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		//expiryErr, _ := err.(*jwt.ValidationError)
		//if expiryErr.Errors != jwt.ValidationErrorExpired {
		//	return "", time.Time{}, err
		//}
		return "", time.Time{}, err
	}

	claims := tokens.Claims.(*Claims)

	issuedAt := time.Now()
	exp := issuedAt.Add(expiry)

	atClaims := jwt.MapClaims{}
	atClaims["iat"] = issuedAt.Unix()
	atClaims["sub"] = claims.Subject
	atClaims["token"] = constant.AccessTokenKey
	atClaims["exp"] = exp.Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)

	newAccessToken, err := at.SignedString([]byte(key))
	if err != nil {
		return "", time.Time{}, err
	}

	return newAccessToken, exp, nil
}
