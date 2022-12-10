package jwthelper

import "github.com/dgrijalva/jwt-go"

var (
	JwtSecretKey      string
	AdminJwtSecretKey string
)

type Claims struct {
	jwt.StandardClaims
	Token   string `json:"token"`
	IsAdmin bool   `json:"isadmin"`
}
