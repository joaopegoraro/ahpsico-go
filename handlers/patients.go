package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/joaopegoraro/ahpsico-go/database/db"
	"github.com/joaopegoraro/ahpsico-go/middlewares"
	"github.com/joaopegoraro/ahpsico-go/server"
)

func HandleShowPatient(s *server.Server) http.HandlerFunc {
	type response struct {
		Uuid        string `json:"uuid"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phoneNumber"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		_, userUuid, err := middlewares.GetAuthDataFromContext(r.Context())
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		patientUuidParam := chi.URLParam(r, "uuid")
		patientUuid, err := uuid.FromString(patientUuidParam)
		if err != nil || patientUuid == uuid.Nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		var patient db.Patient
		if patientUuid == userUuid {
			patient, err = s.Queries.GetPatient(s.Ctx, patientUuid)
		} else {
			patient, err = s.Queries.GetDoctorPatientWithUuid(s.Ctx, db.GetDoctorPatientWithUuidParams{
				DoctorUuid: userUuid,
				Uuid:       patientUuid,
			})
		}

		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		s.RespondOk(w, r, response{
			Uuid:        patient.Uuid.String(),
			Name:        patient.Name,
			PhoneNumber: patient.PhoneNumber,
		})
	}
}

func HandleUpdatePatient(s *server.Server) http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
	}
	type response struct {
		Uuid        string `json:"uuid"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phoneNumber"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		_, userUuid, err := middlewares.GetAuthDataFromContext(r.Context())
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		patientUuidParam := chi.URLParam(r, "uuid")
		patientUuid, err := uuid.FromString(patientUuidParam)
		if err != nil || patientUuid == uuid.Nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		if userUuid != patientUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		var updatedPatient request
		err = s.Decode(w, r, &updatedPatient)
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(updatedPatient.Name) == "" {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		updateParams := db.UpdatePatientParams{}

		patient, err := s.Queries.UpdatePatient(s.Ctx, db.UpdatePatientParams{
			Uuid: patientUuid,
			Name: updateParams.Name,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				s.RespondErrorStatus(w, r, http.StatusNotFound)
				return
			}
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		s.RespondOk(w, r, response{
			Uuid:        patient.Uuid.String(),
			Name:        patient.Name,
			PhoneNumber: patient.PhoneNumber,
		})
	}
}

func HandleListDoctorPatients(s *server.Server) http.HandlerFunc {
	type response struct {
		Uuid        string `json:"uuid"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phoneNumber"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		_, userUuid, err := middlewares.GetAuthDataFromContext(r.Context())
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		doctorUuidQueryParam := r.URL.Query().Get("doctorUuid")
		doctorUuid, err := uuid.FromString(doctorUuidQueryParam)
		if err != nil || doctorUuid == uuid.Nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		if userUuid != doctorUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		fetchedPatients, err := s.Queries.ListDoctorPatients(s.Ctx, doctorUuid)
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
			})
		}

		if len(patients) < 1 {
			s.RespondNoContent(w, r)
			return
		}

		s.RespondOk(w, r, patients)
	}
}
