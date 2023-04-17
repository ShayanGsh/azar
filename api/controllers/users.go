package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Klaushayan/azar/api/pools"
	"github.com/Klaushayan/azar/azar-db"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
    Username    string `json:"username" validate:"required_without=Email,omitempty,min=1,max=100,excludesall=0x20"`
    Email   string `json:"email" validate:"required_without=Username,omitempty,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type UpdateUser struct {
	Username    string `json:"username" validate:"required_without=Email,omitempty,min=1,max=100,excludesall=0x20"`
	Email   string `json:"email" validate:"required_without=Username,omitempty,email"`
	NewEmail    string `json:"new_email" validate:"omitempty,email"`
	OldPassword string `json:"old_password" validate:"required,min=8"`
	NewPassword string `json:"new_password" validate:"min=8"`
}

type UserControllers struct {
	dcp *pools.PGXPool // database connection pool
}

func NewUserControllers(dcp *pools.PGXPool) *UserControllers {
	dcp.SetIdleTimeout(5 * time.Second)
	return &UserControllers{dcp: dcp}
}

func (uc *UserControllers) Login(rw http.ResponseWriter, r *http.Request) {
	var user User
	c, q, err := uc.parseRequest(r, &user)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer uc.dcp.Put(c)

	v, _ := uc.VerifyUser(q, user)
	if v {
		rw.Write([]byte("logged in"))
		return
	}
	rw.Write([]byte("wrong"))
}

func (uc *UserControllers) Register(rw http.ResponseWriter, r *http.Request) {
	var user User
	c, q, err := uc.parseRequest(r, &user)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer uc.dcp.Put(c)

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