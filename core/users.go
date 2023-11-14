package core

import (
	"context"
	"errors"
	"time"

	db "github.com/ShayanGsh/azar/azar-db"
	"github.com/ShayanGsh/azar/core/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

type UpdateUserData struct {
	Username    string `json:"username" validate:"required_without=Email,omitempty,min=1,max=100"`
	NewUsername string `json:"new_username" validate:"omitempty,min=1,max=100"`
	Email   string `json:"email" validate:"required_without=Username,omitempty,email"`
	NewEmail    string `json:"new_email" validate:"omitempty,email"`
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"min=8"`
}

type UserData struct {
    Username    string `json:"username" validate:"required_without=Email,omitempty,min=1,max=100"`
    Email   string `json:"email" validate:"required_without=Username,omitempty,email"`
    Password string `json:"password" validate:"required,min=8"`
}

func AddUser(q *db.Queries, user UserData, context context.Context) error {

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

func AddUserWithHash(q *db.Queries, user UserData, context context.Context) error {

    hashed, err := utils.HashPassword(user.Password)
    if err != nil {
        return err
    }

    user.Password = hashed

    AddUser(q, user, context)
    if err != nil {
        return err
    }

    return nil
}

func UpdateUser(q *db.Queries, updateUser UpdateUserData, context context.Context) error {
    // Get the user by username or email
    user := UserData{
        Username: updateUser.Username,
        Email:    updateUser.Email,
        Password: updateUser.OldPassword,
    }
    existingUser, err := GetUser(q, user, context)
    if err != nil {
        return err
    }

    // Update the user
    UpdateUserFields(updateUser, &existingUser)

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
        err = UpdatePassword(q, updateUser, existingUser, context)
        if err != nil {
            return err
        }
    }

    return nil
}

func GetUser(q *db.Queries, user UserData, context context.Context) (db.User, error) {
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

func UpdateUserFields(updateUser UpdateUserData, existingUser *db.User) {
    if updateUser.NewUsername != "" {
        existingUser.Username = updateUser.NewUsername
    }
    if updateUser.NewEmail != "" {
        existingUser.Email = pgtype.Text{String: updateUser.NewEmail, Valid: true}
    }
    if updateUser.NewPassword != "" {
        // hash the password
        p, err := utils.HashPassword(updateUser.NewPassword)
        if err != nil {
            return
        }
        existingUser.Password = p
    }
}

func UpdatePassword(q *db.Queries, updateUser UpdateUserData, existingUser db.User, context context.Context) error {
    p, err := utils.HashPassword(updateUser.NewPassword)
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

func DeleteUser(q *db.Queries, user UserData, context context.Context) error {
    existingUser, err := GetUser(q, user, context)
    if err != nil {
        return err
    }
    err = q.DeleteUser(context, existingUser.ID)
    if err != nil {
        return err
    }
    return nil
}