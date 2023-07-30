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

	s.Router.Get("/invites", handlers.HandleListInvites(s))
	s.Router.Post("/invites", handlers.HandleCreateInvite(s))
	s.Router.Delete("/invites/{id}", handlers.HandleDeleteInvite(s))
	s.Router.Post("/invites/{id}/accept", handlers.HandleAcceptInvite(s))

	s.Router.Get("/doctors/{uuid}", handlers.HandleShowDoctor(s))
	s.Router.Put("/doctors/{uuid}", handlers.HandleUpdateDoctor(s))
	s.Router.Put("/doctors?patientUuid={patientUuid}", handlers.HandleListPatientDoctors(s))

	s.Router.Get("/patients/{uuid}", handlers.HandleShowPatient(s))
	s.Router.Put("/patients/{uuid}", handlers.HandleUpdatePatient(s))
	s.Router.Put("/patients?doctorUuid={doctorUuid}", handlers.HandleListDoctorPatients(s))
}
