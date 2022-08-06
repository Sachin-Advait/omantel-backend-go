package jwthelper

import "github.com/dgrijalva/jwt-go"

var JwtSecretKey string

type Claims struct {
	jwt.StandardClaims
	Token string
}
