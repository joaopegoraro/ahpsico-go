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
	s.Router.Use(middlewares.Auth(s))

	s.Router.Post("/login", handlers.HandleLoginUser(s))
	s.Router.Post("/register", handlers.HandleRegisterUser(s))

	s.Router.Post("/invites/{id}/accept", handlers.HandleAcceptInvite(s))
	s.Router.Delete("/invites/{id}", handlers.HandleDeleteInvite(s))
	s.Router.Get("/invites", handlers.HandleListInvites(s))
	s.Router.Post("/invites", handlers.HandleCreateInvite(s))

	s.Router.Get("/users/{uuid}", handlers.HandleShowUser(s))
	s.Router.Put("/users/{uuid}", handlers.HandleUpdateUser(s))

	s.Router.Get("/doctors", handlers.HandleListDoctors(s))
	s.Router.Get("/patients", handlers.HandleListPatients(s))

	s.Router.Get("/sessions/{id}", handlers.HandleShowSession(s))
	s.Router.Put("/sessions/{id}", handlers.HandleUpdateSession(s))
	s.Router.Get("/sessions", handlers.HandleListSessions(s))
	s.Router.Post("/sessions", handlers.HandleCreateSession(s))

	s.Router.Put("/assignments/{id}", handlers.HandleUpdateAssignment(s))
	s.Router.Delete("/assignments/{id}", handlers.HandleDeleteAssignment(s))
	s.Router.Get("/assignments", handlers.HandleListAssignments(s))
	s.Router.Post("/assignments", handlers.HandleCreateAssignment(s))

	s.Router.Delete("/advices/{id}", handlers.HandleDeleteAdvice(s))
	s.Router.Get("/advices", handlers.HandleListAdvices(s))
	s.Router.Post("/advices", handlers.HandleCreateAdvice(s))

	s.Router.Delete("/schedule/{id}", handlers.HandleDeleteSchedule(s))
	s.Router.Get("/schedule", handlers.HandleListSchedule(s))
	s.Router.Post("/schedule", handlers.HandleCreateSchedule(s))
}
