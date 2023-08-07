package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/joaopegoraro/ahpsico-go/middlewares"
	"github.com/joaopegoraro/ahpsico-go/server"
)

func HandleListPatients(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		doctorUuidQueryParam := r.URL.Query().Get("doctorUuid")
		if strings.TrimSpace(doctorUuidQueryParam) != "" {
			handleListDoctorPatients(s, doctorUuidQueryParam)(w, r)
			return
		}

		s.RespondErrorStatus(w, r, http.StatusNotFound)
	}
}

func handleListDoctorPatients(s *server.Server, doctorUuidQueryParam string) http.HandlerFunc {
	type response struct {
		Uuid        string `json:"uuid"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phoneNumber"`
		Role        int64  `json:"role"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, userUuid, err := middlewares.GetAuthDataFromContext(ctx)
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		doctorUuid, err := uuid.FromString(doctorUuidQueryParam)
		if err != nil || doctorUuid == uuid.Nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		if userUuid != doctorUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		fetchedPatients, err := s.Queries.ListDoctorPatients(ctx, doctorUuid)
		if err != nil || fetchedPatients == nil {
			if err == sql.ErrNoRows {
				s.RespondNoContent(w, r)
				return
			}
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		patients := []response{}
		for _, patient := range fetchedPatients {
			patients = append(patients, response{
				Uuid:        patient.Uuid.String(),
				Name:        patient.Name,
				PhoneNumber: patient.PhoneNumber,
				Role:        patient.Role,
			})
		}

		if len(patients) < 1 {
			s.RespondNoContent(w, r)
			return
		}

		s.RespondOk(w, r, patients)
	}
}
