package user

import (
	"github.com/ShayanGsh/azar/internal/api"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Controller struct {
	api.Controller
}

func NewController(dcp *pgxpool.Pool, jwt api.JWT) *Controller {
	return &Controller{
		Controller: api.Controller{
			DatabaseConnectionPool: dcp,
			Jwt:                    jwt,
		},
	}
}
