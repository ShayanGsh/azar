package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Klaushayan/azar/api/controllers"
	"github.com/Klaushayan/azar/api/pools"
	"github.com/Klaushayan/azar/azar-db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

type Server struct {
	Router *chi.Mux
	Config *Config
	JWTAuth *jwtauth.JWTAuth
	DB *pools.PGXPool

	// Controllers
	UserControllers *controllers.UserController
}



func NewServer(c *Config) *Server {
	s := &Server {
		Router: chi.NewRouter(),
		Config: c,
		JWTAuth: jwtauth.New("HS256", []byte(c.JWTSecret), nil),
	}

	s.createDBPool()

	s.UserControllers = controllers.NewUserControllers(s.DB)

	return s
}

func (s *Server) setupRoutes() {

	SetupMiddlewares(s.Router)

	s.Router.Group(func(r chi.Router) {
		AuthenticateAccess(r.(*chi.Mux), s.JWTAuth)
	})

	s.Router.Post("/login", s.UserControllers.Login)
	s.Router.Post("/register", s.UserControllers.Register)
}

func (s *Server) createDBPool() {
	connConfig, err := pgx.ParseConfig(s.Config.ToConnString())

	if err != nil {
		panic(err)
	}

	s.DB = pools.NewPool(20, *connConfig)
	s.DB.SetConnTimeout(1 * time.Second)
}

func (s *Server) MigrationCheck() bool{
	conn, err := sql.Open("postgres", s.Config.ToConnString())
	if err != nil {
		panic(err)
	}
	status, err := db.IsMigrated(db.Migration("./"), conn)
	if err != nil {
		panic(err)
	}
	if !status {
		err = db.RunMigration(db.Migration("./"), conn, 0)
		if err != nil {
			panic(err)
		}
		return true
	} else {
		println("Migration already done")
		return true
	}
}

func (s *Server) Start() {
	s.setupRoutes()
	err := http.ListenAndServe(s.Config.Address(), s.Router)
	if err != nil {
		panic(err)
	}
}
