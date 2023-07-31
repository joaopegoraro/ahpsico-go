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

func HandleShowDoctor(s *server.Server) http.HandlerFunc {
	type response struct {
		Uuid           string `json:"uuid"`
		Name           string `json:"name"`
		PhoneNumber    string `json:"phoneNumber"`
		Description    string `json:"description"`
		Crp            string `json:"crp"`
		PixKey         string `json:"pixKey"`
		PaymentDetails string `json:"paymentDetails"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		doctorUuidParam := chi.URLParam(r, "uuid")
		doctorUuid, err := uuid.FromString(doctorUuidParam)
		if err != nil || doctorUuid == uuid.Nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		doctor, err := s.Queries.GetDoctor(ctx, doctorUuid)
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		s.RespondOk(w, r, response{
			Uuid:           doctor.Uuid.String(),
			Name:           doctor.Name,
			PhoneNumber:    doctor.PhoneNumber,
			Description:    doctor.Description,
			Crp:            doctor.Crp,
			PixKey:         doctor.PixKey,
			PaymentDetails: doctor.PaymentDetails,
		})
	}
}

func HandleUpdateDoctor(s *server.Server) http.HandlerFunc {
	type request struct {
		Name           *string `json:"name"`
		Description    *string `json:"description"`
		Crp            *string `json:"crp"`
		PixKey         *string `json:"pixKey"`
		PaymentDetails *string `json:"paymentDetails"`
	}
	type response struct {
		Uuid           string `json:"uuid"`
		Name           string `json:"name"`
		PhoneNumber    string `json:"phoneNumber"`
		Description    string `json:"description"`
		Crp            string `json:"crp"`
		PixKey         string `json:"pixKey"`
		PaymentDetails string `json:"paymentDetails"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, userUuid, err := middlewares.GetAuthDataFromContext(ctx)
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		doctorUuidParam := chi.URLParam(r, "uuid")
		doctorUuid, err := uuid.FromString(doctorUuidParam)
		if err != nil || doctorUuid == uuid.Nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		if userUuid != doctorUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		var updatedDoctor request
		err = s.Decode(w, r, &updatedDoctor)
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		updateParams := db.UpdateDoctorParams{
			Uuid:           doctorUuid,
			Name:           sql.NullString{Valid: false},
			Description:    sql.NullString{Valid: false},
			Crp:            sql.NullString{Valid: false},
			PixKey:         sql.NullString{Valid: false},
			PaymentDetails: sql.NullString{Valid: false},
		}

		if updatedDoctor.Name != nil {
			updateParams.Name = sql.NullString{String: *updatedDoctor.Name, Valid: true}
		}
		if updatedDoctor.Description != nil {
			updateParams.Description = sql.NullString{String: *updatedDoctor.Description, Valid: true}
		}
		if updatedDoctor.Crp != nil {
			updateParams.Crp = sql.NullString{String: *updatedDoctor.Crp, Valid: true}
		}
		if updatedDoctor.PixKey != nil {
			updateParams.PixKey = sql.NullString{String: *updatedDoctor.PixKey, Valid: true}
		}
		if updatedDoctor.PaymentDetails != nil {
			updateParams.PaymentDetails = sql.NullString{String: *updatedDoctor.PaymentDetails, Valid: true}
		}

		doctor, err := s.Queries.UpdateDoctor(ctx, updateParams)
		if err != nil {
			if err == sql.ErrNoRows {
				s.RespondErrorStatus(w, r, http.StatusNotFound)
				return
			}
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		s.RespondOk(w, r, response{
			Uuid:           doctor.Uuid.String(),
			Name:           doctor.Name,
			PhoneNumber:    doctor.PhoneNumber,
			Description:    doctor.Description,
			Crp:            doctor.Crp,
			PixKey:         doctor.PixKey,
			PaymentDetails: doctor.PaymentDetails,
		})
	}
}

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
			})
		}

		if len(doctors) < 1 {
			s.RespondNoContent(w, r)
			return
		}

		s.RespondOk(w, r, doctors)
	}
}
