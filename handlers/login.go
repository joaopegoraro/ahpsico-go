package handlers

import (
	"net/http"

	"github.com/joaopegoraro/ahpsico-go/middlewares"
	"github.com/joaopegoraro/ahpsico-go/server"
)

func HandleLoginUser(s *server.Server) http.HandlerFunc {
	type response struct {
		UserUuid    string `json:"userUuid"`
		UserName    string `json:"userName"`
		PhoneNumber string `json:"phoneNumber"`
		IsDoctor    bool   `json:"isDoctor"`
	}
	var signUpRequiredError = server.Error{
		Type:   "signup_required",
		Detail: "The user is not yet registered in the app",
		Status: http.StatusNotAcceptable,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, userUuid, err := middlewares.GetAuthDataFromContext(ctx)
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		response := response{
			UserUuid:    userUuid.String(),
			UserName:    "",
			PhoneNumber: user.PhoneNumber,
			IsDoctor:    true,
		}

		doctor, err := s.Queries.GetDoctor(ctx, userUuid)
		if err != nil {
			patient, err := s.Queries.GetPatient(ctx, userUuid)
			if err != nil {
				s.RespondError(w, r, signUpRequiredError)
				return
			}
			response.IsDoctor = false
			response.UserName = patient.Name
			s.RespondOk(w, r, response)
			return
		}

		response.UserName = doctor.Name
		s.RespondOk(w, r, response)
	}
}
