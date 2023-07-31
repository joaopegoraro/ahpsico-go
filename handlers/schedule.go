package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/joaopegoraro/ahpsico-go/database/db"
	"github.com/joaopegoraro/ahpsico-go/middlewares"
	"github.com/joaopegoraro/ahpsico-go/server"
	"github.com/joaopegoraro/ahpsico-go/utils"
)

func HandleDeleteSchedule(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, userUuid, err := middlewares.GetAuthDataFromContext(r.Context())
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		scheduleID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || scheduleID < 1 {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		schedule, err := s.Queries.GetSchedule(s.Ctx, int64(scheduleID))
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		if schedule.DoctorUuid != userUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		err = s.Queries.DeleteSchedule(s.Ctx, int64(scheduleID))
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusInternalServerError)
			return
		}

		s.RespondNoContent(w, r)
	}
}

func HandleCreateSchedule(s *server.Server) http.HandlerFunc {
	type request struct {
		DoctorUuid string `json:"doctorUuid"`
		Date       string `json:"date"`
	}
	type response struct {
		ID         int64  `json:"id"`
		DoctorUuid string `json:"doctorUuid"`
		Date       string `json:"date"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		_, userUuid, err := middlewares.GetAuthDataFromContext(r.Context())
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		var createdSchedule request
		err = s.Decode(w, r, &createdSchedule)
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		doctorUuid, err := uuid.FromString(createdSchedule.DoctorUuid)
		if err != nil {
			s.RespondErrorDetail(w, r, "invalid doctor uuid", http.StatusBadRequest)
			return
		}

		if doctorUuid != userUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		parsedDate, err := time.Parse(utils.DateFormat, createdSchedule.Date)
		if err != nil {
			errMessage := fmt.Sprintf("date must be in the following format:  %s", utils.DateFormat)
			s.RespondErrorDetail(w, r, errMessage, http.StatusBadRequest)
			return
		}

		schedule, err := s.Queries.CreateSchedule(s.Ctx, db.CreateScheduleParams{
			DoctorUuid: doctorUuid,
			Date:       parsedDate,
		})
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusInternalServerError)
			return
		}

		s.Respond(w, r, response{
			ID:         schedule.ID,
			DoctorUuid: schedule.DoctorUuid.String(),
			Date:       schedule.Date.Format(utils.DateFormat),
		}, http.StatusCreated)
	}
}

func HandleListDoctorSchedule(s *server.Server) http.HandlerFunc {
	type response struct {
		ID         int64  `json:"id"`
		DoctorUuid string `json:"doctorUuid"`
		Date       string `json:"date"`
		IsSession  bool   `json:"isSession"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		doctorUuidQueryParam := r.URL.Query().Get("doctorUuid")
		doctorUuid, err := uuid.FromString(doctorUuidQueryParam)
		if err != nil || doctorUuid == uuid.Nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		fetchedScheduleList, err := s.Queries.ListDoctorSchedule(s.Ctx, doctorUuid)
		if err != nil || fetchedScheduleList == nil {
			if err == sql.ErrNoRows {
				s.RespondNoContent(w, r)
				return
			}
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		schedule := []response{}
		for _, fetchedSchedule := range fetchedScheduleList {
			schedule = append(schedule, response{
				ID:         fetchedSchedule.ID,
				DoctorUuid: fetchedSchedule.DoctorUuid.String(),
				Date:       fetchedSchedule.Date.Format(utils.DateFormat),
				IsSession:  false,
			})
		}

		fetchedSessions, err := s.Queries.ListDoctorActiveSessions(s.Ctx, doctorUuid)
		if err == nil && len(fetchedSessions) > 0 {
			for _, session := range fetchedSessions {
				schedule = append(schedule, response{
					DoctorUuid: session.DoctorUuid.String(),
					Date:       session.Date.Format(utils.DateFormat),
					IsSession:  true,
				})
			}
		}

		if len(schedule) < 1 {
			s.RespondNoContent(w, r)
			return
		}

		s.RespondOk(w, r, schedule)
	}

}
