package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	db "github.com/Klaushayan/azar/azar-db"
	"github.com/jackc/pgx/v5/pgtype"
)

type UpdateUser struct {
	Username    string `json:"username" validate:"required_without=Email,omitempty,min=1,max=100"`
	NewUsername string `json:"new_username" validate:"omitempty,min=1,max=100"`
	Email   string `json:"email" validate:"required_without=Username,omitempty,email"`
	NewEmail    string `json:"new_email" validate:"omitempty,email"`
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"min=8"`
}

func (uc *UserController) UpdateUserCred(rw http.ResponseWriter, r *http.Request) {
	var u UpdateUser
	c, q, err := uc.parseRequest(r, &u)

	if err != nil {
		ReplyError(rw, err, http.StatusInternalServerError)
		return
	}
	defer c.Release()

	user := User{
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

func (uc *UserController) AddUser(q *db.Queries, user User, context context.Context) error {

    new_user := db.AddUserParams{
        Username: user.Username,
        Email:    pgtype.Text{String: user.Email, Valid: true},
        Password: user.Password,
    }
    err := q.AddUser(context, new_user)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserController) AddUserWithHash(q *db.Queries, user User, context context.Context) error {

    hashed, err := HashPassword(user.Password)
    if err != nil {
        return err
    }

    user.Password = hashed

    uc.AddUser(q, user, context)
    if err != nil {
        return err
    }

    return nil
}

func (uc *UserController) UpdateUser(q *db.Queries, updateUser UpdateUser, context context.Context) error {
    // Get the user by username or email
    user := User{
        Username: updateUser.Username,
        Email:    updateUser.Email,
        Password: updateUser.OldPassword,
    }
    existingUser, err := uc.GetUser(q, user, context)
    if err != nil {
        return err
    }

    // Update the user
    uc.UpdateUserFields(updateUser, &existingUser)

    err = q.UpdateUser(context, db.UpdateUserParams{
        ID:       existingUser.ID,
        Username: existingUser.Username,
        Email:    existingUser.Email,
        UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
    })
    if err != nil {
        return err
    }

    // Update the user's password
    if updateUser.NewPassword != "" {
        err = uc.UpdatePassword(q, updateUser, existingUser, context)
        if err != nil {
            return err
        }
    }

    return nil
}

func (uc *UserController) GetUser(q *db.Queries, user User, context context.Context) (db.User, error) {
    var existingUser db.User
    if user.Username != "" {
        u, err := q.GetUserByUsername(context, user.Username)
        if err != nil {
            return existingUser, err
        }
        existingUser = u
    } else if user.Email != "" {
        u, err := q.GetUserByEmail(context, pgtype.Text{String: user.Email, Valid: true})
        if err != nil {
            return existingUser, err
        }
        existingUser = u
    } else {
        return existingUser, errors.New("username or email must be provided")
    }
    return existingUser, nil
}

func (uc *UserController) UpdateUserFields(updateUser UpdateUser, existingUser *db.User) {
    if updateUser.NewUsername != "" {
        existingUser.Username = updateUser.NewUsername
    }
    if updateUser.NewEmail != "" {
        existingUser.Email = pgtype.Text{String: updateUser.NewEmail, Valid: true}
    }
    if updateUser.NewPassword != "" {
        // hash the password
        p, err := HashPassword(updateUser.NewPassword)
        if err != nil {
            return
        }
        existingUser.Password = p
    }
}

func (uc *UserController) UpdatePassword(q *db.Queries, updateUser UpdateUser, existingUser db.User, context context.Context) error {
    p, err := HashPassword(updateUser.NewPassword)
    if err != nil {
        return err
    }
    err = q.UpdatePassword(context, db.UpdatePasswordParams{
        ID:       existingUser.ID,
        Password: p,
    })
    if err != nil {
        return err
    }
    return nil
}

