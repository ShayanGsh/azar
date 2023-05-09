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

// Server is the main struct for the application. It holds the router, config, and other
// important parts of the application.
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


// NewServer creates a new server instance
// using the provided configuration.
// It initializes a DB pool and
// a new JWTAuth instance.
func NewServer(c *Config) *Server {
	s := &Server {
		Router: chi.NewRouter(),
		Config: c,
		JWTAuth: NewJWT(c.JWTSecret),
	}

	s.createDBPool()

	s.UserControllers = controllers.NewUserController(s.DB, s.JWTAuth)

	return s
}

// SetupRoutes sets up the routes and middleware for the server.
// The routes are grouped by authentication level. The JWTAuth
// middleware is used to verify the token in the header, and
// the AuthenticateAccess middleware is used to verify that the
// user has the correct access level to access the route.
func (s *Server) setupRoutes() {
	SetupMiddlewares(s.Router)
	s.Router.Group(func(r chi.Router) {
		AuthenticateAccess(r.(*chi.Mux), s.JWTAuth.JWTAuth)
	})
	s.Router.Post("/login", s.UserControllers.Login)
	s.Router.Post("/register", s.UserControllers.Register)
}

// createDBPool creates a connection pool to the database.
func (s *Server) createDBPool() {
	connConfig, err := pgxpool.ParseConfig(s.Config.ToConnString())

	if err != nil {
		log.Fatal(err)
	}

	s.DB, err = pgxpool.NewWithConfig(context.Background(), connConfig)

	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) MigrationCheck() bool {
	// Open a connection to the database
	conn, err := sql.Open("postgres", s.Config.ToConnString())
	if err != nil {
		panic(err)
	}

	// Check if the database is migrated
	status, err := db.IsMigrated(db.Migration("./"), conn)
	if err != nil {
		panic(err)
	}

	// If the database is not migrated, run the migration
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

// This function starts a server. It creates a new HTTP server and sets
// the address to the address of the server. It creates a finisher and
// adds the HTTP server to it. It starts listening on the HTTP server
// and sets the finisher to the server. It waits for the finisher to
// finish.
func (s *Server) Start() {
	httpServer := &http.Server{
		Addr: s.Config.Address(),
		Handler: s.Router,
	}
	s.setupRoutes()

	fin := finish.New()
	fin.Add(httpServer)

	s.startListening(httpServer)
	s.Finish = fin
	s.started = true
	fin.Wait()
}

// startListening starts a goroutine to listen for incoming HTTP requests.
// It will also shut down the server if there is an error.
func (s *Server) startListening(httpServer *http.Server) {
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			s.Shutdown()
			s.started = false
			log.Println(err)
		}
	}()
}

// IsRunning returns true if the server is running and false if it is not.
func (s *Server) IsRunning() bool {
	return s.started
}


// Shutdown closes the database connection and
// triggers the Finish channel.
func (s *Server) Shutdown() {
	s.DB.Close()
	s.Finish.Trigger()
	s.started = false
}

// Wait blocks until the server is done.
//
// The server is done when there are no more
// connections.
func (s *Server) Wait() {
	s.Finish.Wait()
}