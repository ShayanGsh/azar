package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Klaushayan/azar/api/pools"
	"github.com/Klaushayan/azar/azar-db"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
    Username    string `json:"username" validate:"required_without=Email,omitempty,min=1,max=100,excludesall=0x20"`
    Email   string `json:"email" validate:"required_without=Username,omitempty,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type UserControllers struct {
	dcp *pools.PGXPool // database connection pool
}

func NewUserControllers(dcp *pools.PGXPool) *UserControllers {
	dcp.SetIdleTimeout(5 * time.Second)
	return &UserControllers{dcp: dcp}
}

func (uc *UserControllers) Login(rw http.ResponseWriter, r *http.Request) {
	c, err := uc.dcp.Get()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer uc.dcp.Put(c)

	q := db.New(c)
	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	err = validator.New().Struct(user)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	v, _ := uc.VerifyUser(q, user)
	if v {
		rw.Write([]byte("logged in"))
		return
	}
	rw.Write([]byte("wrong"))
}

func (uc *UserControllers) Register(rw http.ResponseWriter, r *http.Request) {
	c, err := uc.dcp.Get()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer uc.dcp.Put(c)

	q := db.New(c)
	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	err = validator.New().Struct(user)
    if err != nil {
        rw.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
    }
	if err := q.AddUser(r.Context(), db.AddUserParams{
		Username: user.Username,
		Email:   pgtype.Text{String: user.Email},
		Password: user.Password,
	}); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func (uc *UserControllers) VerifyUser(queries *db.Queries, user User) (bool, error) {
	var err error
	ctx := context.Background()
	if user.Username != "" {
		err = queries.VerifyUser(ctx, db.VerifyUserParams{Username: user.Username, Password: user.Password})
	}
	if user.Email != "" {
		err = queries.VerifyUserByEmail(ctx, db.VerifyUserByEmailParams{Email: pgtype.Text{String: user.Email}, Password: user.Password})
	}
	if err != nil {
		return false, err
	}
	return true, nil
}