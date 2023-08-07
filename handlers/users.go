package handlers

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/joaopegoraro/ahpsico-go/database/db"
	"github.com/joaopegoraro/ahpsico-go/middlewares"
	"github.com/joaopegoraro/ahpsico-go/server"
)

const (
	patientRole = iota
	doctorRole
)

const (
	firstUserRole = patientRole
	lastUserRole  = doctorRole
)

func HandleShowUser(s *server.Server) http.HandlerFunc {
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

		userUuidParam := chi.URLParam(r, "uuid")
		userUuid, err := uuid.FromString(userUuidParam)
		if err != nil || userUuid == uuid.Nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		user, err := s.Queries.GetUser(ctx, userUuid)
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		s.RespondOk(w, r, response{
			Uuid:           user.Uuid.String(),
			Name:           user.Name,
			PhoneNumber:    user.PhoneNumber,
			Description:    user.Description,
			Crp:            user.Crp,
			PixKey:         user.PixKey,
			PaymentDetails: user.PaymentDetails,
			Role:           user.Role,
		})
	}
}

func HandleUpdateUser(s *server.Server) http.HandlerFunc {
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
		Role           int64  `json:"role"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, userUuid, err := middlewares.GetAuthDataFromContext(ctx)
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		userUuidParam := chi.URLParam(r, "uuid")
		userUuidFromParam, err := uuid.FromString(userUuidParam)
		if err != nil || userUuidFromParam == uuid.Nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		if userUuid != userUuidFromParam {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		var updatedUser request
		err = s.Decode(w, r, &updatedUser)
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		updateParams := db.UpdateUserParams{
			Uuid:           userUuid,
			Name:           sql.NullString{Valid: false},
			Description:    sql.NullString{Valid: false},
			Crp:            sql.NullString{Valid: false},
			PixKey:         sql.NullString{Valid: false},
			PaymentDetails: sql.NullString{Valid: false},
		}

		if updatedUser.Name != nil {
			updateParams.Name = sql.NullString{String: *updatedUser.Name, Valid: true}
		}
		if updatedUser.Description != nil {
			updateParams.Description = sql.NullString{String: *updatedUser.Description, Valid: true}
		}
		if updatedUser.Crp != nil {
			updateParams.Crp = sql.NullString{String: *updatedUser.Crp, Valid: true}
		}
		if updatedUser.PixKey != nil {
			updateParams.PixKey = sql.NullString{String: *updatedUser.PixKey, Valid: true}
		}
		if updatedUser.PaymentDetails != nil {
			updateParams.PaymentDetails = sql.NullString{String: *updatedUser.PaymentDetails, Valid: true}
		}

		user, err := s.Queries.UpdateUser(ctx, updateParams)
		if err != nil {
			if err == sql.ErrNoRows {
				s.RespondErrorStatus(w, r, http.StatusNotFound)
				return
			}
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		s.RespondOk(w, r, response{
			Uuid:           user.Uuid.String(),
			Name:           user.Name,
			PhoneNumber:    user.PhoneNumber,
			Description:    user.Description,
			Crp:            user.Crp,
			PixKey:         user.PixKey,
			PaymentDetails: user.PaymentDetails,
			Role:           user.Role,
		})
	}
}
