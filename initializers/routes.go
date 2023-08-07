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

	s.Router.Post("/verification-code", handlers.HandleSendVerificationCode(s))
	s.Router.Post("/login", handlers.HandleLoginUser(s))

	s.Router.Route("/signup", func(r chi.Router) {
		r.Use(middlewares.Auth(s, true))
		r.Post("/", handlers.HandleCreateUser(s))
	})

	s.Router.Route("/", func(r chi.Router) {
		r.Use(middlewares.Auth(s, false))

		r.Post("/invites/{id}/accept", handlers.HandleAcceptInvite(s))
		r.Delete("/invites/{id}", handlers.HandleDeleteInvite(s))
		r.Get("/invites", handlers.HandleListInvites(s))
		r.Post("/invites", handlers.HandleCreateInvite(s))

		r.Get("/users/{uuid}", handlers.HandleShowUser(s))
		r.Put("/users/{uuid}", handlers.HandleUpdateUser(s))

		r.Get("/doctors", handlers.HandleListDoctors(s))
		r.Get("/patients", handlers.HandleListPatients(s))

		r.Get("/sessions/{id}", handlers.HandleShowSession(s))
		r.Put("/sessions/{id}", handlers.HandleUpdateSession(s))
		r.Get("/sessions", handlers.HandleListSessions(s))
		r.Post("/sessions", handlers.HandleCreateSession(s))

		r.Put("/assignments/{id}", handlers.HandleUpdateAssignment(s))
		r.Delete("/assignments/{id}", handlers.HandleDeleteAssignment(s))
		r.Get("/assignments", handlers.HandleListAssignments(s))
		r.Post("/assignments", handlers.HandleCreateAssignment(s))

		r.Delete("/advices/{id}", handlers.HandleDeleteAdvice(s))
		r.Get("/advices", handlers.HandleListAdvices(s))
		r.Post("/advices", handlers.HandleCreateAdvice(s))

		r.Delete("/schedule/{id}", handlers.HandleDeleteSchedule(s))
		r.Get("/schedule", handlers.HandleListSchedule(s))
		r.Post("/schedule", handlers.HandleCreateSchedule(s))
	})
}
