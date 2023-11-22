package user

import (
	"github.com/ShayanGsh/azar/internal/rest"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Controller struct {
	rest.Controller
}

func NewController(dcp *pgxpool.Pool, jwt rest.JWT) *Controller {
	return &Controller{
		Controller: rest.Controller{
			DatabaseConnectionPool: dcp,
			Jwt:                    jwt,
		},
	}
}
