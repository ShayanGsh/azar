package user

import (
	"errors"
	"github.com/ShayanGsh/azar/internal/rest"
	"github.com/ShayanGsh/azar/internal/model"
	"net/http"
)

func (uc *Controller) UpdateUserCred(rw http.ResponseWriter, r *http.Request) {
	var u model.UpdateUserData
	c, q, err := uc.ParseRequest(r, &u)

	if err != nil {
		rest.ReplyError(rw, err, http.StatusInternalServerError)
		return
	}
	defer c.Release()

	user := model.UserData{
		Username: u.Username,
		Email:    u.Email,
		Password: u.OldPassword,
	}

	v, err := uc.VerifyUser(q, user)
	if err != nil {
		rest.ReplyError(rw, err, http.StatusInternalServerError)
		return
	}
	if v {
		err = model.UpdateUser(q, u, r.Context())
		if err != nil {
			rest.ReplyError(rw, err, http.StatusInternalServerError)
			return
		}
		rest.ReplySuccess(rw, "user updated successfully", http.StatusOK)
	} else {
		rest.ReplyError(rw, errors.New("invalid credentials"), http.StatusUnauthorized)
	}
}
