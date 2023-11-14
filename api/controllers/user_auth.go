package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/ShayanGsh/azar/azar-db"
	"github.com/ShayanGsh/azar/core"
	"github.com/ShayanGsh/azar/core/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserController struct {
	Controller
}

func NewUserController(dcp *pgxpool.Pool, jwt JWT) *UserController {
	return &UserController{
		Controller: Controller{
			DatabaseConnectionPool: dcp,
			jwt: jwt,
		},
	}
}

func (uc *UserController) Login(rw http.ResponseWriter, r *http.Request) {
	var user core.UserData
	c, q, err := uc.parseRequest(r, &user)

	if err != nil {
		ReplyError(rw, err, http.StatusInternalServerError)
		return
	}
	defer c.Release()

	v, err := uc.VerifyUser(q, user)
	if err != nil {
		ReplyError(rw, err, http.StatusInternalServerError)
		return
	}
	if v {
		token, _, err := uc.jwt.Encode(user.Username, user.Username) //TODO: Add some sort of id (not the db id) to the token
		if err != nil {
			ReplyError(rw, err, http.StatusInternalServerError)
			return
		}
		ReplySuccess(rw, token, http.StatusOK)
		return
	}
	ReplyError(rw, errors.New("invalid credentials"), http.StatusUnauthorized)
}

func (uc *UserController) Register(rw http.ResponseWriter, r *http.Request) {
	var user core.UserData
	c, q, err := uc.parseRequest(r, &user)

	if err != nil {
		ReplyError(rw, err, http.StatusInternalServerError)
		return
	}
	defer c.Release()

	passwordStrength := utils.GetPasswordStrength(user.Password)
	// TODO: Add password strength check using the minimum requirements in the config file
	if passwordStrength == utils.VeryWeak {
		ReplyError(rw, errors.New("password is too weak"), http.StatusForbidden)
		return
	}
	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		ReplyError(rw, errors.New("password hashing failed"), http.StatusInternalServerError)
		return
	}

	if err := q.AddUser(r.Context(), db.AddUserParams{
		Username: user.Username,
		Email:   pgtype.Text{String: user.Email},
		Password: hashedPassword,
	}); err != nil {
		dbe := uc.parseDBError(err)
		log.Println(dbe)
		switch dbe.Code {
		case "23505":
			ReplyError(rw, errors.New("username or email already exists"), http.StatusConflict)
			return
		default:
			ReplyError(rw, errors.New(dbe.Message), http.StatusInternalServerError)
			return
		}
	}
	ReplySuccess(rw, "success", http.StatusCreated)
}

func (uc *UserController) VerifyUser(queries *db.Queries, user core.UserData) (bool, error) {
	if user.Email != "" {
		email := pgtype.Text{String: user.Email}
		u, err := queries.GetUserByEmail(context.Background(), email)
		if err != nil {
			return false, err
		}
		return utils.CheckPasswordHash(user.Password, u.Password), nil
	}
	u, err := queries.GetUserByUsername(context.Background(), user.Username)
	if err != nil {
		return false, err
	}
	return utils.CheckPasswordHash(user.Password, u.Password), nil
}