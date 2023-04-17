package controllers

import (
	"encoding/json"
	"net/http"

	db "github.com/Klaushayan/azar/azar-db"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
)

func (uc *UserControllers) parseRequest(r *http.Request, body interface{}) (*pgx.Conn, *db.Queries, error) {

	c, err := uc.dcp.Get()
	if err != nil {
		return nil, nil, err
	}

	defer uc.dcp.Put(c)

	q := db.New(c)

	if r.Body == nil {
		return nil, nil, err
	}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, nil, err
	}
	err = validator.New().Struct(body)
	if err != nil {
		return nil, nil, err
	}
	return c, q, nil
}
