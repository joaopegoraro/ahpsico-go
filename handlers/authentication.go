package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/joaopegoraro/ahpsico-go/middlewares"
	"github.com/joaopegoraro/ahpsico-go/server"
	"github.com/joaopegoraro/ahpsico-go/utils"

	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

func HandleSendVerificationCode(s *server.Server) http.HandlerFunc {
	type request struct {
		PhoneNumber string `json:"phoneNumber"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO INIT
		s.RespondNoContent(w, r)
		return
		// TODO END
		var phoneRequest request
		err := s.Decode(w, r, &phoneRequest)
		if err != nil || strings.TrimSpace(phoneRequest.PhoneNumber) == "" {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		params := &openapi.CreateVerificationParams{}
		params.SetTo(phoneRequest.PhoneNumber)
		params.SetChannel("sms")

		_, err = s.Twilio.VerifyV2.CreateVerification(os.Getenv("TWILIO_VERIFY_SERVICE_SID"), params)
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		s.RespondNoContent(w, r)
	}
}

func HandleLoginUser(s *server.Server) http.HandlerFunc {
	type request struct {
		PhoneNumber string `json:"phoneNumber"`
		Code        string `json:"code"`
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
	var signUpRequiredError = server.Error{
		Type:   "signup_required",
		Detail: "The user is not yet registered in the app",
		Status: http.StatusNotAcceptable,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var loginRequest request
		err := s.Decode(w, r, &loginRequest)
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(loginRequest.PhoneNumber) == "" || strings.TrimSpace(loginRequest.Code) == "" {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		// TODO UNCCOMENT
		//	params := &openapi.CreateVerificationCheckParams{}
		//	params.SetTo(loginRequest.PhoneNumber)
		//	params.SetCode(loginRequest.Code)

		//	resp, err := s.Twilio.VerifyV2.CreateVerificationCheck(os.Getenv("TWILIO_VERIFY_SERVICE_SID"), params)

		//	if err != nil || *resp.Status != "approved" {
		//		s.RespondErrorStatus(w, r, http.StatusBadRequest)
		//		return
		//	}

		user, err := s.Queries.GetUserByPhoneNumber(ctx, loginRequest.PhoneNumber)
		if err != nil {
			newUuid, err := uuid.NewV4()
			if err != nil {
				s.RespondErrorStatus(w, r, http.StatusInternalServerError)
				return
			}

			token, err := utils.GenerateJWT(newUuid.String(), loginRequest.PhoneNumber, middlewares.TemporaryUserRole)
			if err != nil {
				s.RespondErrorStatus(w, r, http.StatusInternalServerError)
				return
			}

			middlewares.SetTokenHeader(w, token)
			s.RespondError(w, r, signUpRequiredError)
			return
		}

		token, err := utils.GenerateJWT(user.Uuid.String(), user.PhoneNumber, user.Role)
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusInternalServerError)
			return
		}

		response := response{
			Uuid:           user.Uuid.String(),
			Name:           user.Name,
			PhoneNumber:    user.PhoneNumber,
			Description:    user.Description,
			Crp:            user.Crp,
			PixKey:         user.PixKey,
			PaymentDetails: user.PaymentDetails,
			Role:           user.Role,
		}

		middlewares.SetTokenHeader(w, token)
		s.RespondOk(w, r, response)
	}
}
