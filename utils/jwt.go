package utils

import (
	"demo-store/common"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("my_secret_key")

var TokenExpirationInMinutes int = 30
var BearerTokenHeader = "Bearer "

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Tokenizer interface {
	CreateToken(username string) (string, error)
	GetUsernameFromToken(tokenString string) (string, error)
}

type JwtTokenizer struct {
	Tracer Tracer
}

func NewJwtTokenizer(tracer Tracer) *JwtTokenizer {
	return &JwtTokenizer{Tracer: tracer}
}

func (j *JwtTokenizer) CreateToken(username string) (string, error) {

	expirationTime := time.Now().Add(time.Minute * time.Duration(TokenExpirationInMinutes))
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "UserJWTService",
		}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		j.Tracer.LogError("Error creating the token: ", err)
		return "", common.ErrorCreatingJwtToken
	}

	return tokenString, nil
}

func (j *JwtTokenizer) GetUsernameFromToken(tokenString string) (string, error) {

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			j.Tracer.LogError("Error signature invalid ", err)
			return "", common.ErrorAuthorizationFailed
		}
		j.Tracer.LogError("Error processing JWT token ", err)
		return "", common.ErrorAuthorizationFailed
	}

	if !token.Valid {
		j.Tracer.LogError("Error invalid token ", err)
		return "", common.ErrorAuthorizationFailed
	}

	return claims.Username, nil
}
