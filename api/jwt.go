package api

import (
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type JWT struct {
	*jwtauth.JWTAuth
	DefaultExpiration time.Duration
}

func NewJWT(secret string) *JWT {
	return &JWT {
		JWTAuth: jwtauth.New("HS256", []byte(secret), nil),
		DefaultExpiration: time.Hour * 1, // default expiration time is 1 hour
	}
}

func (j *JWT) Encode(userID, username string, expiration ...time.Duration) (string, jwt.Token, error) {

	var exp time.Duration
	if len(expiration) > 0 {
		exp = expiration[0]
	} else {
		exp = j.DefaultExpiration
	}

	t, tokenString, err := j.JWTAuth.Encode(map[string]interface{} {
		"id": userID,
		"username": username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(exp).Unix(),
	})

	if err != nil {
		return "", nil, err
	}

	return tokenString, t, nil
}