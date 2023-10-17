package mpb

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	Claims
	jwt.RegisteredClaims
}

type JWTAdminClaims struct {
	AdminClaims
	jwt.RegisteredClaims
}
