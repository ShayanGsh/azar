package user

import (
	"context"
	"errors"
	"github.com/ShayanGsh/azar/internal/rest"
	"github.com/ShayanGsh/azar/internal/model"
	"github.com/ShayanGsh/azar/internal/utils"
	"log"
	"net/http"

	"github.com/ShayanGsh/azar/azar-db"
	"github.com/jackc/pgx/v5/pgtype"
)

func (uc *Controller) Login(rw http.ResponseWriter, r *http.Request) {
	var user model.UserData
	c, q, err := uc.ParseRequest(r, &user)

	if err != nil {
		rest.ReplyError(rw, err, http.StatusInternalServerError)
		return
	}
	defer c.Release()

	v, err := uc.VerifyUser(q, user)
	if err != nil {
		rest.ReplyError(rw, err, http.StatusInternalServerError)
		return
	}
	if v {
		token, _, err := uc.Jwt.Encode(user.Username, user.Username) //TODO: Add some sort of id (not the db id) to the token
		if err != nil {
			rest.ReplyError(rw, err, http.StatusInternalServerError)
			return
		}
		rest.ReplySuccess(rw, token, http.StatusOK)
		return
	}
	rest.ReplyError(rw, errors.New("invalid credentials"), http.StatusUnauthorized)
}

func (uc *Controller) Register(rw http.ResponseWriter, r *http.Request) {
	var user model.UserData
	c, q, err := uc.ParseRequest(r, &user)

	if err != nil {
		rest.ReplyError(rw, err, http.StatusInternalServerError)
		return
	}
	defer c.Release()

	passwordStrength := utils.GetPasswordStrength(user.Password)
	// TODO: Add password strength check using the minimum requirements in the config file
	if passwordStrength == utils.VeryWeak {
		rest.ReplyError(rw, errors.New("password is too weak"), http.StatusForbidden)
		return
	}
	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		rest.ReplyError(rw, errors.New("password hashing failed"), http.StatusInternalServerError)
		return
	}

	if err := q.AddUser(r.Context(), db.AddUserParams{
		Username: user.Username,
		Email:    pgtype.Text{String: user.Email},
		Password: hashedPassword,
	}); err != nil {
		dbe := uc.ParseDBError(err)
		log.Println(dbe)
		switch dbe.Code {
		case "23505":
			rest.ReplyError(rw, errors.New("username or email already exists"), http.StatusConflict)
			return
		default:
			rest.ReplyError(rw, errors.New(dbe.Message), http.StatusInternalServerError)
			return
		}
	}
	rest.ReplySuccess(rw, "success", http.StatusCreated)
}

func (uc *Controller) VerifyUser(queries *db.Queries, user model.UserData) (bool, error) {
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
