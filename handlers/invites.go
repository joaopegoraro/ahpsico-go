package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/joaopegoraro/ahpsico-go/database/db"
	"github.com/joaopegoraro/ahpsico-go/middlewares"
	"github.com/joaopegoraro/ahpsico-go/server"
)

func HandleListInvites(s *server.Server) http.HandlerFunc {
	type doctor struct {
		Uuid        string `json:"uuid"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	type response struct {
		ID          int64  `json:"id"`
		PhoneNumber string `json:"phoneNumber"`
		PatientUuid string `json:"patientUuid"`
		Doctor      doctor `json:"doctor"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, userUuid, err := middlewares.GetAuthDataFromContext(ctx)
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		invites := []response{}

		fetchedInvites, err := s.Queries.ListDoctorInvites(ctx, userUuid)
		if err != nil || len(fetchedInvites) == 0 {
			fetchedPatientInvites, err := s.Queries.ListPatientInvites(ctx, userUuid)
			if err != nil || len(fetchedPatientInvites) < 1 {
				s.RespondNoContent(w, r)
				return
			}

			for _, invite := range fetchedPatientInvites {
				fetchedInvites = append(fetchedInvites, db.ListDoctorInvitesRow(invite))
			}
		}

		for _, invite := range fetchedInvites {
			invites = append(invites, response{
				ID:          invite.InviteID,
				PhoneNumber: invite.InvitePhoneNumber,
				PatientUuid: invite.InvitePatientUuid.String(),
				Doctor: doctor{
					Uuid:        invite.DoctorUuid.String(),
					Name:        invite.DoctorName,
					Description: invite.DoctorDescription,
				},
			})
		}

		if len(invites) < 1 {
			s.RespondNoContent(w, r)
			return
		}

		s.RespondOk(w, r, invites)
	}
}

func HandleCreateInvite(s *server.Server) http.HandlerFunc {
	type doctor struct {
		Uuid        string `json:"uuid"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	type request struct {
		PhoneNumber string `json:"phoneNumber"`
	}
	type response struct {
		ID          int64  `json:"id"`
		PhoneNumber string `json:"phoneNumber"`
		PatientUuid string `json:"patientUuid"`
		Doctor      doctor `json:"doctor"`
	}
	var inviteAlreadySentError = server.Error{
		Type:   "invite_already_sent",
		Detail: "There already exists an invite from this doctor to this patient",
		Status: http.StatusConflict,
	}
	var patientNotRegisteredError = server.Error{
		Type:   "patient_not_registered",
		Detail: "There are no patients registered with this phone number yet",
		Status: http.StatusNotFound,
	}
	var patientAlreadyWithDoctorError = server.Error{
		Type:   "patient_already_with_doctor",
		Detail: "You can't send the invite, the patient is already with the doctor",
		Status: http.StatusConflict,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, userUuid, err := middlewares.GetAuthDataFromContext(ctx)
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		user, err := s.Queries.GetUserByRole(ctx, db.GetUserByRoleParams{
			Uuid: userUuid,
			Role: middlewares.DoctorRole,
		})
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		var newInvite request
		err = s.Decode(w, r, &newInvite)
		if err != nil || strings.TrimSpace(newInvite.PhoneNumber) == "" {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		_, err = s.Queries.GetDoctorInviteByPhoneNumber(ctx, db.GetDoctorInviteByPhoneNumberParams{
			DoctorUuid:  userUuid,
			PhoneNumber: newInvite.PhoneNumber,
		})
		if err == nil {
			s.RespondError(w, r, inviteAlreadySentError)
			return
		}

		patient, err := s.Queries.GetUserByPhoneNumber(ctx, newInvite.PhoneNumber)
		if err != nil {
			s.RespondError(w, r, patientNotRegisteredError)
			return
		}

		patients, _ := s.Queries.ListDoctorPatientsByPhoneNumber(ctx, db.ListDoctorPatientsByPhoneNumberParams{
			DoctorUuid:  userUuid,
			PhoneNumber: newInvite.PhoneNumber,
		})
		if len(patients) > 0 {
			s.RespondError(w, r, patientAlreadyWithDoctorError)
			return
		}

		invite, err := s.Queries.CreateInvite(ctx, db.CreateInviteParams{
			PhoneNumber: newInvite.PhoneNumber,
			DoctorUuid:  userUuid,
			PatientUuid: patient.Uuid,
		})
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusInternalServerError)
		}

		s.Respond(w, r, response{
			ID:          invite.ID,
			PhoneNumber: invite.PhoneNumber,
			PatientUuid: invite.PatientUuid.String(),
			Doctor: doctor{
				Uuid:        invite.DoctorUuid.String(),
				Name:        user.Name,
				Description: user.Description,
			},
		}, http.StatusCreated)
	}
}

func HandleDeleteInvite(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, userUuid, err := middlewares.GetAuthDataFromContext(ctx)
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		inviteId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		invite, err := s.Queries.GetInvite(ctx, int64(inviteId))
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		if invite.DoctorUuid != userUuid && invite.PatientUuid != userUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		err = s.Queries.DeleteInvite(ctx, invite.ID)
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusInternalServerError)
			return
		}

		s.RespondNoContent(w, r)
	}
}

func HandleAcceptInvite(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, userUuid, err := middlewares.GetAuthDataFromContext(ctx)
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		inviteId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		invite, err := s.Queries.GetInvite(ctx, int64(inviteId))
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		if invite.PatientUuid != userUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		err = s.Queries.AddPatientDoctor(ctx, db.AddPatientDoctorParams{
			DoctorUuid:  invite.DoctorUuid,
			PatientUuid: invite.PatientUuid,
		})
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusInternalServerError)
			return
		}

		err = s.Queries.DeleteInvite(ctx, invite.ID)
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusInternalServerError)
			return
		}

		s.RespondNoContent(w, r)
	}
}
