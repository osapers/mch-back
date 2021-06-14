package httpapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/osapers/mch-back/internal/constant"
)

const ttl = 24 * time.Hour * 365 // 1 year

func (s *Server) tokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		tokenStr := c.GetHeader("Authorization")

		if len(tokenStr) > 7 {
			tokenStr = tokenStr[7:]
		}

		userID, err := parseToken(tokenStr, s.cfg.secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, jsonResp(nil, err))
		}

		ctx := context.WithValue(c.Request.Context(), constant.UserIDCtxKey, userID)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

type customClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func parseToken(tokenStr string, secret []byte) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return "", fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*customClaims)

	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token: %v", token.Claims)
	}

	return claims.UserID, nil
}

func newToken(userID string, secret []byte) string {
	claims := customClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, _ := token.SignedString(secret)

	return ss
}
