package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var appSecret = []byte("ecom-app-secret-key")

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}


func GenerateToken(userID, role string)(string, error){
	claims := Claims{
		UserID: userID,
		Role: role,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 *time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(appSecret)
}


func GenerateRefreshToken(userID string)(string, time.Time, error){
	exp := time.Now().Add(7* 24 * time.Hour)

	claims := jwt.RegisteredClaims{
		Subject: userID,
		ExpiresAt: jwt.NewNumericDate(exp),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}

	token :=jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(appSecret)
	return tokenStr, exp, err
}



func ParseToken(tokenStr string) (*Claims, error) {

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (interface{}, error) {

			// Ensure signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return appSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
