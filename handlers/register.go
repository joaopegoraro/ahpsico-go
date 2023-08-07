package handlers

import (
	"net/http"
	"strings"

	"github.com/joaopegoraro/ahpsico-go/database/db"
	"github.com/joaopegoraro/ahpsico-go/middlewares"
	"github.com/joaopegoraro/ahpsico-go/server"
)

func HandleRegisterUser(s *server.Server) http.HandlerFunc {
	type request struct {
		UserName string `json:"name"`
		IsDoctor *bool  `json:"isDoctor"`
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
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}
		if newUser.IsDoctor == nil || strings.TrimSpace(newUser.UserName) == "" {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		_, err = s.Queries.GetUser(ctx, userUuid)
		if err == nil {
			s.RespondError(w, r, userAlreadyRegisteredError)
			return
		}

		var role int64
		if *newUser.IsDoctor {
			role = doctorRole
		} else {
			role = patientRole
		}

		createdUser, err := s.Queries.CreateUser(ctx, db.CreateUserParams{
			Uuid:        userUuid,
			Name:        newUser.UserName,
			PhoneNumber: user.PhoneNumber,
			Role:        role,
		})

		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
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

		s.Respond(w, r, response, http.StatusCreated)
	}
}
