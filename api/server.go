package api

import (
	"context"
	"database/sql"
	"net/http"
	"log"

	"github.com/Klaushayan/azar/api/controllers"
	"github.com/Klaushayan/azar/azar-db"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pseidemann/finish"
)

type Server struct {
	Router *chi.Mux
	Config *Config
	JWTAuth *JWT
	DB *pgxpool.Pool
	Finish *finish.Finisher

	started bool

	// Controllers
	UserControllers *controllers.UserController
}


func NewServer(c *Config) *Server {
	s := &Server {
		Router: chi.NewRouter(),
		Config: c,
		JWTAuth: NewJWT(c.JWTSecret),
	}

	s.createDBPool()

	s.UserControllers = controllers.NewUserControllers(s.DB)

	return s
}

func (s *Server) setupRoutes() {

	SetupMiddlewares(s.Router)

	s.Router.Group(func(r chi.Router) {
		AuthenticateAccess(r.(*chi.Mux), s.JWTAuth.JWTAuth)
	})

	s.Router.Post("/login", s.UserControllers.Login)
	s.Router.Post("/register", s.UserControllers.Register)
}

func (s *Server) createDBPool() {
	connConfig, err := pgxpool.ParseConfig(s.Config.ToConnString())

	if err != nil {
		panic(err)
	}

	s.DB, err = pgxpool.NewWithConfig(context.Background(), connConfig)

	if err != nil {
		panic(err)
	}
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
	httpServer := &http.Server{
		Addr: s.Config.Address(),
		Handler: s.Router,
	}
	s.setupRoutes()

	fin := finish.New()
	fin.Add(httpServer)

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			s.Shutdown()
			s.started = false
			log.Println(err)
		}
	}()
	s.Finish = fin
	s.started = true
	fin.Wait()
}

func (s *Server) IsRunning() bool {
	return s.started
}

func (s *Server) Shutdown() {
	s.DB.Close()
	s.Finish.Trigger()
	s.started = false
}

func (s *Server) Wait() {
	s.Finish.Wait()
}