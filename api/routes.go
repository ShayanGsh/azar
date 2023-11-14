package api

import "github.com/go-chi/chi/v5"

var AUTH_API_PATH = "/api/auth"

func SetRoutes(mux *chi.Mux, s *APIServer) {
	mux.Route(AUTH_API_PATH, func(r chi.Router) {
		r.Post("/login", s.UserControllers.Login)
		r.Post("/register", s.UserControllers.Register)
	})

	mux.Group(func(r chi.Router) {
		AuthenticateAccess(r.(*chi.Mux), s.JWTAuth.JWTAuth)

		r.Route("/api", func(r chi.Router) {
			r.Patch("/users", s.UserControllers.UpdateUserCred)
			// r.Get("/users", s.UserControllers.GetUsers)
		})
	})
}