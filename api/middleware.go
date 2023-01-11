package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
)

func SetupMiddlewares(router *chi.Mux) *chi.Mux {
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           10, // 300 is the maximum value not ignored by any of major browsers
		}))

	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	return router
}

func AuthenticateAccess(router *chi.Mux, jwt *jwtauth.JWTAuth) *chi.Mux {
	router.Use(jwtauth.Verifier(jwt))
	router.Use(jwtauth.Authenticator)
	return router
}