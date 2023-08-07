package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/joaopegoraro/ahpsico-go/middlewares"
	"github.com/joaopegoraro/ahpsico-go/server"
)

func HandleListDoctors(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		patientUuidQueryParam := r.URL.Query().Get("patientUuid")
		if strings.TrimSpace(patientUuidQueryParam) != "" {
			handleListPatientDoctors(s, patientUuidQueryParam)(w, r)
			return
		}

		s.RespondErrorStatus(w, r, http.StatusNotFound)
	}
}

func handleListPatientDoctors(s *server.Server, patientUuidQueryParam string) http.HandlerFunc {
	type response struct {
		Uuid           string `json:"uuid"`
		Name           string `json:"name"`
		PhoneNumber    string `json:"phoneNumber"`
		Description    string `json:"description"`
		Crp            string `json:"crp"`
		PixKey         string `json:"pixKey"`
		PaymentDetails string `json:"paymentDetails"`
		Role           int64  `json:"role"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, userUuid, err := middlewares.GetAuthDataFromContext(ctx)
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		patientUuid, err := uuid.FromString(patientUuidQueryParam)
		if err != nil || patientUuid == uuid.Nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		if userUuid != patientUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		fetchedDoctors, err := s.Queries.ListPatientDoctors(ctx, patientUuid)
		if err != nil || fetchedDoctors == nil {
			if err == sql.ErrNoRows {
				s.RespondNoContent(w, r)
				return
			}
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		doctors := []response{}
		for _, doctor := range fetchedDoctors {
			doctors = append(doctors, response{
				Uuid:           doctor.Uuid.String(),
				Name:           doctor.Name,
				PhoneNumber:    doctor.PhoneNumber,
				Description:    doctor.Description,
				Crp:            doctor.Crp,
				PixKey:         doctor.PixKey,
				PaymentDetails: doctor.PaymentDetails,
				Role:           doctor.Role,
			})
		}

		if len(doctors) < 1 {
			s.RespondNoContent(w, r)
			return
		}

		s.RespondOk(w, r, doctors)
	}
}
