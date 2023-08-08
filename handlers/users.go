package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/joaopegoraro/ahpsico-go/database/db"
	"github.com/joaopegoraro/ahpsico-go/middlewares"
	"github.com/joaopegoraro/ahpsico-go/server"
	"github.com/joaopegoraro/ahpsico-go/utils"
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

func HandleCreateUser(s *server.Server) http.HandlerFunc {
	type request struct {
		UserName string `json:"name"`
		Role     int64  `json:"role"`
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
	var userAlreadyRegisteredError = server.Error{
		Type:   "user_already_registered",
		Detail: "The user is already registered in the app",
		Status: http.StatusNotAcceptable,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, userUuid, err := middlewares.GetAuthDataFromContext(ctx)
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		var newUser request
		err = s.Decode(w, r, &newUser)
		if err != nil {
			s.RespondError(w, r, server.Error{
				Detail: err.Error(),
				Status: http.StatusBadRequest,
			})
			return
		}
		if strings.TrimSpace(newUser.UserName) == "" {
			s.RespondError(w, r, server.Error{
				Detail: "Name cannot be blank",
				Status: http.StatusBadRequest,
			})
			return
		}
		if newUser.Role < middlewares.FirstUserRole || newUser.Role > middlewares.LastUserRole {
			s.RespondError(w, r, server.Error{
				Detail: fmt.Sprintf("Role must be a valid role. Role provided: %d", newUser.Role),
				Status: http.StatusBadRequest,
			})
			return
		}

		_, err = s.Queries.GetUser(ctx, userUuid)
		if err == nil {
			s.RespondError(w, r, userAlreadyRegisteredError)
			return
		}

		createdUser, err := s.Queries.CreateUser(ctx, db.CreateUserParams{
			Uuid:        userUuid,
			Name:        newUser.UserName,
			PhoneNumber: user.PhoneNumber,
			Role:        newUser.Role,
		})

		if err != nil {
			s.RespondError(w, r, server.Error{
				Detail: err.Error(),
				Status: http.StatusBadRequest,
			})
			return
		}

		response := response{
			Uuid:           createdUser.Uuid.String(),
			Name:           createdUser.Name,
			PhoneNumber:    createdUser.PhoneNumber,
			Description:    createdUser.Description,
			Crp:            createdUser.Crp,
			PixKey:         createdUser.PixKey,
			PaymentDetails: createdUser.PaymentDetails,
			Role:           createdUser.Role,
		}

		tokenTeste, erro := utils.GenerateJWT(response.Uuid, response.PhoneNumber, response.Role)
		if erro != nil {
			s.RespondErrorStatus(w, r, http.StatusInternalServerError)
			return
		}

		middlewares.SetTokenHeader(w, tokenTeste)
		s.Respond(w, r, response, http.StatusCreated)
	}
}
