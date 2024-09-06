package models

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type JWTHandler struct {
	config *viper.Viper
}

func NewJWTHandler(config *viper.Viper) *JWTHandler {
	return &JWTHandler{
		config: config,
	}
}

func (jh *JWTHandler) GenerateJWT(id string) (string, error) {
	duration := jh.config.GetDuration("jwt.token_hour_lifespan")
	secret := jh.config.GetString("jwt.api_secret")
	claims := jwt.RegisteredClaims{
		Subject:   id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (jh *JWTHandler) GetJWTFromHeader(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""

}

func (jh *JWTHandler) GetIDFromJWT(tokenString string) (string, error) {
	secret := jh.config.GetString("jwt.api_secret")
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", err
	}
	// check if token is not expired
	if !token.Valid {
		return "", err
	}
	return claims.Subject, nil
}

func (jh *JWTHandler) GetIDFromToken(token string) string {
	id, _ := jh.GetIDFromJWT(token)
	return id
}

func (jh *JWTHandler) GetIDFromHeader(c *gin.Context) string {
	tokenString := jh.GetJWTFromHeader(c)
	id, _ := jh.GetIDFromJWT(tokenString)
	return id
}

func (jh *JWTHandler) TokenValid(c *gin.Context) error {
	tokenString := jh.GetJWTFromHeader(c)
	secret := jh.config.GetString("jwt.api_secret")
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return nil
	}
	return err
}
