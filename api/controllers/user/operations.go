package user

import (
	"errors"
	"github.com/ShayanGsh/azar/internal/api"
	"github.com/ShayanGsh/azar/internal/models"
	"net/http"
)

func (uc *Controller) UpdateUserCred(rw http.ResponseWriter, r *http.Request) {
	var u models.UpdateUserData
	c, q, err := uc.ParseRequest(r, &u)

	if err != nil {
		api.ReplyError(rw, err, http.StatusInternalServerError)
		return
	}
	defer c.Release()

	user := models.UserData{
		Username: u.Username,
		Email:    u.Email,
		Password: u.OldPassword,
	}

	v, err := uc.VerifyUser(q, user)
	if err != nil {
		api.ReplyError(rw, err, http.StatusInternalServerError)
		return
	}
	if v {
		err = models.UpdateUser(q, u, r.Context())
		if err != nil {
			api.ReplyError(rw, err, http.StatusInternalServerError)
			return
		}
		api.ReplySuccess(rw, "user updated successfully", http.StatusOK)
	} else {
		api.ReplyError(rw, errors.New("invalid credentials"), http.StatusUnauthorized)
	}
}
