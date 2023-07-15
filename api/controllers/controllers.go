package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	db "github.com/ShayanGsh/azar/azar-db"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type JWT interface {
	Encode(userID, username string, expiration ...time.Duration) (string, jwt.Token, error)
}

type Controller struct {
	DatabaseConnectionPool *pgxpool.Pool
	jwt JWT
}

type DBError struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
	Detail  string `json:"Detail"`
}

func (ctrl *Controller) parseRequest(r *http.Request, body any) (*pgxpool.Conn, *db.Queries, error) {

	c, err := ctrl.DatabaseConnectionPool.Acquire(r.Context())
	if err != nil {
		return nil, nil, errors.New("internal server error")
	}

	q := db.New(c)

	if r.Body == nil {
		return c, q, nil
	}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, nil, errors.New("invalid request body")
	}
	err = validator.New().Struct(body)
	if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
			// TODO: to clean up later
            if err.Tag() == "min" {
				return nil, nil, errors.New("password must be at least 8 characters long")
            }
			if err.Tag() == "email" {
				return nil, nil, errors.New("invalid email address")
			}
			if err.Tag() == "required" {
				return nil, nil, errors.New("missing required field")
			}
        }
		return nil, nil, errors.New("invalid request body")
	}
	return c, q, nil
}

func (ctrl *Controller) parseDBError(err error) DBError {
	if err == nil {
		return DBError{}
	}
	var dbErr DBError

	pgErr, ok := err.(*pgconn.PgError)
	dbErr.Code = pgErr.Code
	dbErr.Message = pgErr.Message
	dbErr.Detail = pgErr.Detail

	if !ok {
		return DBError{
			Code:    "0",
			Message: "internal server error",
			Detail:  "internal server error",
		}
	}

	return dbErr
}