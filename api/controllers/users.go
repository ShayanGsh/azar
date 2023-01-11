package controllers

import (
	"net/http"
	"github.com/Klaushayan/azar/api/db"
	"time"
)

type UserControllers struct {
	db *db.Pool
}

func NewUserControllers(db *db.Pool) *UserControllers {
	db.SetIdleTimeout(5 * time.Second)
	return &UserControllers{db: db}
}

func (uc *UserControllers) Login(rw http.ResponseWriter, r *http.Request) {

}