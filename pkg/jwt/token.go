package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type (
	TokenOptions struct {
		AccessSecret string
		AccessExpire int64
		Fields       map[string]interface{}
	}

	Token struct {
		AccessToken  string `json:"access_token"`
		AccessExpire int64  `json:"access_expire"`
	}

	Claims struct {
		UserID int64 `json:"user_id"`
		jwt.RegisteredClaims
	}
)

func BuildTokens(opt TokenOptions) (Token, error) {
	// 定义载荷
	ExpireTime := time.Now().Add(time.Duration(opt.AccessExpire) * time.Second)
	claims := &Claims{
		opt.Fields["userId"].(int64),
		jwt.RegisteredClaims{
			Issuer:    "ccmo-gozero-blog",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(opt.AccessExpire) * time.Second)),
		},
	}
	// 签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(opt.AccessSecret))
	if err != nil {
		return Token{}, err
	}

	return Token{
		tokenString,
		ExpireTime.Unix(),
	}, nil
}

func ParseToken(tokenString, accessSecret string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
