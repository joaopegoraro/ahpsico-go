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
		UserUuid    string `json:"userUuid"`
		UserName    string `json:"userName"`
		PhoneNumber string `json:"phoneNumber"`
		IsDoctor    bool   `json:"isDoctor"`
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

		_, err = s.Queries.GetDoctor(ctx, userUuid)
		if err == nil {
			s.RespondError(w, r, userAlreadyRegisteredError)
			return
		}
		_, err = s.Queries.GetPatient(ctx, userUuid)
		if err == nil {
			s.RespondError(w, r, userAlreadyRegisteredError)
			return
		}

		if *newUser.IsDoctor {
			_, err = s.Queries.CreateDoctor(ctx, db.CreateDoctorParams{
				Uuid:        userUuid,
				Name:        newUser.UserName,
				PhoneNumber: user.PhoneNumber,
			})
		} else {
			_, err = s.Queries.CreatePatient(ctx, db.CreatePatientParams{
				Uuid:        userUuid,
				Name:        newUser.UserName,
				PhoneNumber: user.PhoneNumber,
			})
		}

		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		response := response{
			UserUuid:    userUuid.String(),
			UserName:    newUser.UserName,
			PhoneNumber: user.PhoneNumber,
			IsDoctor:    *newUser.IsDoctor,
		}

		s.Respond(w, r, response, http.StatusCreated)
	}
}
