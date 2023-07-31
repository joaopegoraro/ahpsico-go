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
	s.Router.Get("/doctors?patientUuid={patientUuid}", handlers.HandleListPatientDoctors(s))

	s.Router.Get("/patients/{uuid}", handlers.HandleShowPatient(s))
	s.Router.Put("/patients/{uuid}", handlers.HandleUpdatePatient(s))
	s.Router.Get("/patients?doctorUuid={doctorUuid}", handlers.HandleListDoctorPatients(s))

	s.Router.Get("/sessions/{id}", handlers.HandleShowSession(s))
	s.Router.Put("/sessions/{id}", handlers.HandleUpdateSession(s))
	s.Router.Get("/sessions?doctorUuid={doctorUuid}&date={date}", handlers.HandleListDoctorSessions(s))
	s.Router.Get("/sessions?patientUuid={patientUuid}&upcoming={upcoming}", handlers.HandleListPatientSessions(s))
	s.Router.Post("/sessions", handlers.HandleCreateSession(s))

	s.Router.Put("/assignments/{id}", handlers.HandleUpdateAssignment(s))
	s.Router.Delete("/assignments/{id}", handlers.HandleDeleteAssignment(s))
	s.Router.Get("/assignments?patientUuid={patientUuid}&pending={pending}", handlers.HandleListPatientAssignments(s))
	s.Router.Post("/assignments", handlers.HandleCreateAssignment(s))

	s.Router.Delete("/advices/{id}", handlers.HandleDeleteAdvice(s))
	s.Router.Get("/advices?patientUuid={patientUuid}", handlers.HandleListPatientAdvices(s))
	s.Router.Get("/advices?doctorUuid={doctorUuid}", handlers.HandleListDoctorAdvices(s))
	s.Router.Post("/advices", handlers.HandleCreateAdvice(s))

	s.Router.Delete("/schedule/{id}", handlers.HandleDeleteSchedule(s))
	s.Router.Get("/schedule?doctorUuid={doctorUuid}", handlers.HandleListDoctorSchedule(s))
	s.Router.Post("/schedule", handlers.HandleCreateSchedule(s))
}
