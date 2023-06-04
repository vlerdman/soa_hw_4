package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

var JWT_KEY = []byte("8Zz5tw0Ionm3XPZZfN0NOml3z9FMfmpgXwovR9fp6ryDIoGRM8EPHAB6iHsc0fb")

type Claims struct {
	UserId uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

func GetToken(userId uuid.UUID) string {
	claims := &Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(JWT_KEY)
	return tokenString
}

func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWT_KEY, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func FetchToken(w http.ResponseWriter, r *http.Request) (*Claims, error) {
	token := r.Header.Get("Authorization")
	claims, err := ParseToken(token)
	if err != nil {
		return nil, err
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("expired token")
	}
	return claims, nil
}
