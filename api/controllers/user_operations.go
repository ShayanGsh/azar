package controllers

import (
	"errors"
	"net/http"
	"github.com/ShayanGsh/azar/core"
)

func (uc *UserController) UpdateUserCred(rw http.ResponseWriter, r *http.Request) {
	var u core.UpdateUserData
	c, q, err := uc.parseRequest(r, &u)

	if err != nil {
		ReplyError(rw, err, http.StatusInternalServerError)
		return
	}
	defer c.Release()

	user := core.UserData{
		Username: u.Username,
		Email: u.Email,
		Password: u.OldPassword,
	}

	v, err := uc.VerifyUser(q, user)
	if err != nil {
		ReplyError(rw, err, http.StatusInternalServerError)
		return
	}
	if v {
		ReplySuccess(rw, "updated user", http.StatusOK)
		return
	}
	ReplyError(rw, errors.New("invalid credentials"), http.StatusUnauthorized)
}