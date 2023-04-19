package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	db "github.com/Klaushayan/azar/azar-db"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Controller struct {
	dcp *pgxpool.Pool
}

type DBError struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
	Detail  string `json:"Detail"`
}

func (ctrl *Controller) parseRequest(r *http.Request, body interface{}) (*pgxpool.Conn, *db.Queries, error) {

	c, err := ctrl.dcp.Acquire(r.Context())
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