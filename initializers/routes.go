package initializers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joaopegoraro/ahpsico-go/handlers"
	"github.com/joaopegoraro/ahpsico-go/middlewares"
	"github.com/joaopegoraro/ahpsico-go/server"
)

func InitializeRoutes(s *server.Server) {
	s.Router = chi.NewRouter()

	s.Router.Use(middleware.Logger)
	s.Router.Use(middlewares.Security(s))

	s.Router.Post("login", handlers.HandleLoginUser(s))
	s.Router.Post("register", handlers.HandleRegisterUser(s))
}
