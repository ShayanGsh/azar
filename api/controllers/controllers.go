package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"time"

	db "github.com/Klaushayan/azar/azar-db"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/crypto/bcrypt"
)

type JWT interface {
	Encode(userID, username string, expiration ...time.Duration) (string, jwt.Token, error)
}

type Controller struct {
	dcp *pgxpool.Pool
	jwt JWT
}

type DBError struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
	Detail  string `json:"Detail"`
}

func (ctrl *Controller) parseRequest(r *http.Request, body any) (*pgxpool.Conn, *db.Queries, error) {

	c, err := ctrl.dcp.Acquire(r.Context())
	if err != nil {
		return nil, nil, errors.New("internal server error")
	}

	q := db.New(c)

	if r.Body == nil {
		return c, q, nil
	}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, nil, errors.New("invalid request body")
	}
	err = validator.New().Struct(body)
	if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
			// TODO: to clean up later
            if err.Tag() == "min" {
				return nil, nil, errors.New("password must be at least 8 characters long")
            }
			if err.Tag() == "email" {
				return nil, nil, errors.New("invalid email address")
			}
			if err.Tag() == "required" {
				return nil, nil, errors.New("missing required field")
			}
        }
		return nil, nil, errors.New("invalid request body")
	}
	return c, q, nil
}

func (ctrl *Controller) parseDBError(err error) DBError {
	if err == nil {
		return DBError{}
	}
	var dbErr DBError

	pgErr, ok := err.(*pgconn.PgError)
	dbErr.Code = pgErr.Code
	dbErr.Message = pgErr.Message
	dbErr.Detail = pgErr.Detail

	if !ok {
		return DBError{
			Code:    "0",
			Message: "internal server error",
			Detail:  "internal server error",
		}
	}

	return dbErr
}

type PasswordStrength string

const (
	VeryWeak   PasswordStrength = "very weak"
	Weak       PasswordStrength = "weak"
	Normal      PasswordStrength = "normal"
	Strong      PasswordStrength = "strong"
	Unbeatable  PasswordStrength = "unbeatable"
)

func GetPasswordStrength(password string) PasswordStrength {
	length := len(password)

	if length < 8 {
		return VeryWeak
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()]`).MatchString(password)

	score := 0
	if hasLower {
		score++
	}
	if hasUpper {
		score++
	}
	if hasDigit {
		score++
	}
	if hasSpecial {
		score++
	}
	if length >= 12 {
		score++
	}

	switch score {
	case 1, 2:
		return Weak
	case 3:
		return Normal
	case 4:
		return Strong
	default:
		return Unbeatable
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}