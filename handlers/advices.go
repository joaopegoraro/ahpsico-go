package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/joaopegoraro/ahpsico-go/database/db"
	"github.com/joaopegoraro/ahpsico-go/middlewares"
	"github.com/joaopegoraro/ahpsico-go/server"
	"github.com/joaopegoraro/ahpsico-go/utils"
)

func HandleDeleteAdvice(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, userUuid, err := middlewares.GetAuthDataFromContext(r.Context())
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		adviceID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || adviceID < 1 {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		advice, err := s.Queries.GetAdvice(s.Ctx, int64(adviceID))
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		if advice.DoctorUuid != userUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		err = s.Queries.DeleteAdvice(s.Ctx, int64(adviceID))
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusInternalServerError)
			return
		}

		s.RespondNoContent(w, r)
	}
}

func HandleCreateAdvice(s *server.Server) http.HandlerFunc {
	type request struct {
		DoctorUuid   string   `json:"doctorUuid"`
		Message      string   `json:"message"`
		PatientUuids []string `json:"patientUuids"`
	}
	type response struct {
		ID           int64    `json:"id"`
		DoctorUuid   string   `json:"doctorUuid"`
		Message      string   `json:"message"`
		PatientUuids []string `json:"patientUuids"`
		CreatedAt    string   `json:"createdAt"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		_, userUuid, err := middlewares.GetAuthDataFromContext(r.Context())
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		var createdAdvice request
		err = s.Decode(w, r, &createdAdvice)
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		doctorUuid, err := uuid.FromString(createdAdvice.DoctorUuid)
		if err != nil {
			s.RespondErrorDetail(w, r, "invalid doctor uuid", http.StatusBadRequest)
			return
		}

		if doctorUuid != userUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		if strings.TrimSpace(createdAdvice.Message) == "" {
			s.RespondErrorDetail(w, r, "message must not be empty", http.StatusBadRequest)
			return
		}

		patientUuids := []uuid.UUID{}
		for _, patientUuid := range createdAdvice.PatientUuids {
			uuid, err := uuid.FromString(patientUuid)
			if err != nil {
				errMessage := fmt.Sprintf("invalid patient uuid %s", uuid.String())
				s.RespondErrorDetail(w, r, errMessage, http.StatusBadRequest)
				return
			}
			patientUuids = append(patientUuids, uuid)
		}

		if len(patientUuids) < 1 {
			s.RespondErrorDetail(w, r, "patient uuids must not be empty", http.StatusBadRequest)
			return
		}

		advice, err := s.Queries.CreateAdvice(s.Ctx, db.CreateAdviceParams{
			DoctorUuid: doctorUuid,
			Message:    createdAdvice.Message,
		})
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusInternalServerError)
			return
		}

		for _, patientUuid := range patientUuids {
			err = s.Queries.CreateAdviceWithPatient(s.Ctx, db.CreateAdviceWithPatientParams{
				AdviceID:    advice.ID,
				PatientUuid: patientUuid,
			})
			if err != nil {
				s.RespondErrorStatus(w, r, http.StatusInternalServerError)
				return
			}
		}

		s.Respond(w, r, response{
			ID:           advice.ID,
			DoctorUuid:   advice.DoctorUuid.String(),
			Message:      advice.Message,
			PatientUuids: createdAdvice.PatientUuids,
		}, http.StatusCreated)
	}
}

func HandleListAdvices(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		doctorUuidQueryParam := r.URL.Query().Get("doctorUuid")
		if strings.TrimSpace(doctorUuidQueryParam) != "" {
			handleListDoctorAdvices(s, doctorUuidQueryParam)(w, r)
			return
		}

		patientUuidQueryParam := r.URL.Query().Get("patientUuid")
		if strings.TrimSpace(patientUuidQueryParam) != "" {
			handleListPatientAdvices(s, patientUuidQueryParam)(w, r)
			return
		}

		s.RespondErrorStatus(w, r, http.StatusNotFound)
	}
}

func handleListDoctorAdvices(s *server.Server, doctorUuidQueryParam string) http.HandlerFunc {
	type doctor struct {
		Uuid string `json:"uuid"`
		Name string `json:"name"`
	}
	type response struct {
		ID           int64    `json:"id"`
		Doctor       doctor   `json:"doctor"`
		Message      string   `json:"message"`
		PatientUuids []string `json:"patientUuids"`
		CreatedAt    string   `json:"createdAt"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		_, userUuid, err := middlewares.GetAuthDataFromContext(r.Context())
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		doctorUuid, err := uuid.FromString(doctorUuidQueryParam)
		if err != nil || doctorUuid == uuid.Nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		if userUuid != doctorUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		fetchedAdvices, err := s.Queries.ListDoctorAdvices(s.Ctx, doctorUuid)
		if err != nil || fetchedAdvices == nil {
			if err == sql.ErrNoRows {
				s.RespondNoContent(w, r)
				return
			}
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		advices := []response{}
		patientUuidsMap := map[int64][]string{}
		for _, fetchedAdvice := range fetchedAdvices {
			advicePatientUuids := patientUuidsMap[fetchedAdvice.AdviceID]
			patientUuidsMap[fetchedAdvice.AdviceID] = append(advicePatientUuids, fetchedAdvice.PatientUuid.String())

			advices = append(advices, response{
				ID:      fetchedAdvice.AdviceID,
				Message: fetchedAdvice.AdviceMessage,
				Doctor: doctor{
					Uuid: fetchedAdvice.DoctorUuid.String(),
					Name: fetchedAdvice.DoctorName,
				},
				CreatedAt: fetchedAdvice.AdviceCreatedAt.Format(utils.DateFormat),
			})
		}

		for index, advice := range advices {
			advice.PatientUuids = patientUuidsMap[advice.ID]
			advices[index] = advice
		}

		if len(advices) < 1 {
			s.RespondNoContent(w, r)
			return
		}

		s.RespondOk(w, r, advices)
	}

}

func handleListPatientAdvices(s *server.Server, patientUuidQueryParam string) http.HandlerFunc {
	type doctor struct {
		Uuid string `json:"uuid"`
		Name string `json:"name"`
	}
	type response struct {
		ID          int64  `json:"id"`
		Doctor      doctor `json:"doctor"`
		Message     string `json:"message"`
		PatientUuid string `json:"patientUuid"`
		CreatedAt   string `json:"createdAt"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		_, userUuid, err := middlewares.GetAuthDataFromContext(r.Context())
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		patientUuid, err := uuid.FromString(patientUuidQueryParam)
		if err != nil || patientUuid == uuid.Nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		var fetchedAdvices []db.ListPatientAdvicesRow
		if userUuid == patientUuid {
			fetchedAdvices, err = s.Queries.ListPatientAdvices(s.Ctx, patientUuid)
		} else {
			var list []db.ListDoctorPatientAdvicesRow
			list, err = s.Queries.ListDoctorPatientAdvices(s.Ctx, db.ListDoctorPatientAdvicesParams{
				PatientUuid: patientUuid,
				DoctorUuid:  userUuid,
			})
			for _, advice := range list {
				fetchedAdvices = append(fetchedAdvices, db.ListPatientAdvicesRow(advice))
			}
		}

		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		advices := []response{}
		for _, fetchedAdvice := range fetchedAdvices {
			advices = append(advices, response{
				ID:          fetchedAdvice.AdviceID,
				Message:     fetchedAdvice.AdviceMessage,
				PatientUuid: patientUuid.String(),
				Doctor: doctor{
					Uuid: fetchedAdvice.DoctorUuid.String(),
					Name: fetchedAdvice.DoctorName,
				},
				CreatedAt: fetchedAdvice.AdviceCreatedAt.Format(utils.DateFormat),
			})
		}

		if len(advices) < 1 {
			s.RespondNoContent(w, r)
			return
		}

		s.RespondOk(w, r, advices)
	}

}
