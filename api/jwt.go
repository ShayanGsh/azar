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
	t, tokenString, err := j.JWTAuth.Encode(map[string]interface{} {
		"id": userID,
		"username": username,
	})

	if err != nil {
		return "", nil, err
	}

	if len(expiration) > 0 {
		t.Expiration().Add(expiration[0])
	} else {
		t.Expiration().Add(time.Hour * 24)
	}

	return tokenString, t, nil
}