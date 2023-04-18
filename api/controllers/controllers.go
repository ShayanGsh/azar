package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Klaushayan/azar/api/pools"
	db "github.com/Klaushayan/azar/azar-db"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
)

type Controller struct {
	dcp *pools.PGXPool // database connection pool
}

func (ctrl *Controller) parseRequest(r *http.Request, body interface{}) (*pgx.Conn, *db.Queries, error) {

	c, err := ctrl.dcp.Get()
	if err != nil {
		return nil, nil, errors.New("internal server error")
	}

	defer ctrl.dcp.Put(c)

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