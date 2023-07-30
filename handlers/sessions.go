package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/joaopegoraro/ahpsico-go/database/db"
	"github.com/joaopegoraro/ahpsico-go/middlewares"
	"github.com/joaopegoraro/ahpsico-go/server"
	"github.com/joaopegoraro/ahpsico-go/utils"
)

const (
	notConfirmedStatus int64 = iota
	confirmedStatus
	canceledStatus
	concludedStatus
)
const (
	firstStatus = notConfirmedStatus
	_           = confirmedStatus
	_           = canceledStatus
	lastStatus  = concludedStatus
)

const (
	individualType int64 = iota
	monthlyType
)
const (
	firstType = individualType
	lastType  = monthlyType
)

func HandleShowSession(s *server.Server) http.HandlerFunc {
	type doctor struct {
		Uuid        string `json:"uuid"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	type patient struct {
		Uuid        string `json:"uuid"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phoneNumber"`
	}
	type response struct {
		ID         int64   `json:"id"`
		Doctor     doctor  `json:"doctor"`
		Patient    patient `json:"patient"`
		GroupIndex int64   `json:"groupIndex"`
		Status     int64   `json:"status"`
		Type       int64   `json:"type"`
		Date       string  `json:"date"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		_, userUuid, err := middlewares.GetAuthDataFromContext(r.Context())
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		sessionId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || sessionId < 1 {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		session, err := s.Queries.GetSessionWithParticipants(s.Ctx, int64(sessionId))
		if err != nil || sessionId < 1 {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		if session.DoctorUuid != userUuid && session.PatientUuid != userUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		s.RespondOk(w, r, response{
			ID:         session.SessionID,
			GroupIndex: session.SessionGroupIndex,
			Status:     session.SessionStatus,
			Type:       session.SessionType,
			Date:       session.SessionDate.Format(utils.DateFormat),
			Doctor: doctor{
				Uuid:        session.DoctorUuid.String(),
				Name:        session.DoctorName,
				Description: session.DoctorDescription,
			},
			Patient: patient{
				Uuid:        session.PatientUuid.String(),
				Name:        session.PatientName,
				PhoneNumber: session.PatientPhoneNumber,
			},
		})
	}
}

func HandleCreateSession(s *server.Server) http.HandlerFunc {
	type request struct {
		DoctorUuid  string `json:"doctorUuid"`
		PatientUuid string `json:"patientUuid"`
		GroupIndex  int64  `json:"groupIndex"`
		Status      int64  `json:"status"`
		Type        int64  `json:"type"`
		Date        string `json:"date"`
	}
	type response struct {
		ID          int64  `json:"id"`
		DoctorUuid  string `json:"doctorUuid"`
		PatientUuid string `json:"patientUuid"`
		GroupIndex  int64  `json:"groupIndex"`
		Status      int64  `json:"status"`
		Type        int64  `json:"type"`
		Date        string `json:"date"`
	}
	var sessionAlreadyBookedError = server.Error{
		Type:   "session_already_booked",
		Detail: "There already is a session booked at this time.",
		Status: http.StatusConflict,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		_, userUuid, err := middlewares.GetAuthDataFromContext(r.Context())
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		var createdSession request
		err = s.Decode(w, r, &createdSession)
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		doctorUuid, err := uuid.FromString(createdSession.DoctorUuid)
		if err != nil {
			s.RespondErrorDetail(w, r, "invalid doctor uuid", http.StatusBadRequest)
			return
		}

		patientUuid, err := uuid.FromString(createdSession.PatientUuid)
		if err != nil {
			s.RespondErrorDetail(w, r, "invalid patient uuid", http.StatusBadRequest)
			return
		}

		if doctorUuid != userUuid && patientUuid != userUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		if createdSession.Status < firstStatus || createdSession.Status > lastStatus {
			errMessage := fmt.Sprintf("status must be between %d and %d", firstStatus, lastStatus)
			s.RespondErrorDetail(w, r, errMessage, http.StatusBadRequest)
			return
		}
		if createdSession.Type < firstType || createdSession.Type > lastType {
			errMessage := fmt.Sprintf("type must be between %d and %d", firstType, lastType)
			s.RespondErrorDetail(w, r, errMessage, http.StatusBadRequest)
			return
		}

		parsedDate, err := time.Parse(utils.DateFormat, createdSession.Date)
		if err != nil {
			errMessage := fmt.Sprintf("date must be in the following format:  %s", utils.DateFormat)
			s.RespondErrorDetail(w, r, errMessage, http.StatusBadRequest)
			return
		}

		_, err = s.Queries.GetDoctorSessionByExactDate(s.Ctx, db.GetDoctorSessionByExactDateParams{
			DoctorUuid: doctorUuid,
			Date:       parsedDate,
		})
		if err == nil {
			s.RespondError(w, r, sessionAlreadyBookedError)
			return
		}

		session, err := s.Queries.CreateSession(s.Ctx, db.CreateSessionParams{
			PatientUuid: patientUuid,
			DoctorUuid:  doctorUuid,
			Date:        parsedDate,
			GroupIndex:  createdSession.GroupIndex,
			Status:      createdSession.Status,
			Type:        createdSession.Type,
		})
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		s.Respond(w, r, response{
			ID:          session.ID,
			DoctorUuid:  session.DoctorUuid.String(),
			PatientUuid: session.PatientUuid.String(),
			GroupIndex:  session.GroupIndex,
			Status:      session.Status,
			Type:        session.Type,
			Date:        session.Date.Format(utils.DateFormat),
		}, http.StatusCreated)
	}
}

func HandleUpdateSession(s *server.Server) http.HandlerFunc {
	type request struct {
		Status *int64  `json:"status"`
		Date   *string `json:"date"`
	}
	type response struct {
		ID          int64  `json:"id"`
		DoctorUuid  string `json:"doctorUuid"`
		PatientUuid string `json:"patientUuid"`
		GroupIndex  int64  `json:"groupIndex"`
		Status      int64  `json:"status"`
		Type        int64  `json:"type"`
		Date        string `json:"date"`
	}
	var sessionAlreadyBookedError = server.Error{
		Type:   "session_already_booked",
		Detail: "There already is a session booked at this time.",
		Status: http.StatusConflict,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		_, userUuid, err := middlewares.GetAuthDataFromContext(r.Context())
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		sessionId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || sessionId < 1 {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		var updatedSession request
		err = s.Decode(w, r, &updatedSession)
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		savedSession, err := s.Queries.GetSession(s.Ctx, int64(sessionId))
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		if savedSession.DoctorUuid != userUuid && savedSession.PatientUuid != userUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		var updateSessionParams = db.UpdateSessionParams{
			ID:     int64(sessionId),
			Status: sql.NullInt64{Valid: false},
			Date:   sql.NullTime{Valid: false},
		}

		if updatedSession.Status != nil {
			if *updatedSession.Status < firstStatus || *updatedSession.Status > lastStatus {
				errMessage := fmt.Sprintf("status must be between %d and %d", firstStatus, lastStatus)
				s.RespondErrorDetail(w, r, errMessage, http.StatusBadRequest)
				return
			}
			updateSessionParams.Status = sql.NullInt64{Int64: *updatedSession.Status, Valid: true}
		}

		if updatedSession.Date != nil {
			parsedDate, err := time.Parse(utils.DateFormat, *updatedSession.Date)
			if err != nil {
				errMessage := fmt.Sprintf("date must be in the following format:  %s", utils.DateFormat)
				s.RespondErrorDetail(w, r, errMessage, http.StatusBadRequest)
				return
			}

			_, err = s.Queries.GetDoctorSessionByExactDate(s.Ctx, db.GetDoctorSessionByExactDateParams{
				DoctorUuid: savedSession.DoctorUuid,
				Date:       parsedDate,
			})
			if err == nil {
				s.RespondError(w, r, sessionAlreadyBookedError)
				return
			}

			updateSessionParams.Date = sql.NullTime{Time: parsedDate, Valid: true}
		}

		session, err := s.Queries.UpdateSession(s.Ctx, updateSessionParams)
		if err != nil {
			if err == sql.ErrNoRows {
				s.RespondErrorStatus(w, r, http.StatusNotFound)
				return
			}
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		s.RespondOk(w, r, response{
			ID:          session.ID,
			DoctorUuid:  session.DoctorUuid.String(),
			PatientUuid: session.PatientUuid.String(),
			GroupIndex:  session.GroupIndex,
			Status:      session.Status,
			Type:        session.Type,
			Date:        session.Date.Format(utils.DateFormat),
		})
	}
}

func HandleListDoctorSessions(s *server.Server) http.HandlerFunc {
	type patient struct {
		Uuid        string `json:"uuid"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phoneNumber"`
	}
	type response struct {
		ID         int64   `json:"id"`
		Patient    patient `json:"patient"`
		GroupIndex int64   `json:"groupIndex"`
		Status     int64   `json:"status"`
		Type       int64   `json:"type"`
		Date       string  `json:"date"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		_, userUuid, err := middlewares.GetAuthDataFromContext(r.Context())
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		doctorUuidQueryParam := r.URL.Query().Get("doctorUuid")
		doctorUuid, err := uuid.FromString(doctorUuidQueryParam)
		if err != nil || doctorUuid == uuid.Nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		if userUuid != doctorUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		dateParam := r.URL.Query().Get("date")

		sessions := []response{}

		if strings.TrimSpace(dateParam) == "" {
			fetchedSessions, err := s.Queries.ListDoctorSessions(s.Ctx, doctorUuid)
			if err != nil || fetchedSessions == nil {
				if err == sql.ErrNoRows {
					s.RespondNoContent(w, r)
					return
				}
				s.RespondErrorStatus(w, r, http.StatusNotFound)
				return
			}

			for _, session := range fetchedSessions {
				sessions = append(sessions, response{
					ID:         session.SessionID,
					GroupIndex: session.SessionGroupIndex,
					Status:     session.SessionStatus,
					Type:       session.SessionType,
					Date:       session.SessionDate.Format(utils.DateFormat),
					Patient: patient{
						Uuid:        session.PatientUuid.String(),
						Name:        session.PatientName,
						PhoneNumber: session.PatientPhoneNumber,
					},
				})
			}
		} else {
			parsedDate, err := time.Parse(utils.DateFormat, dateParam)
			if err != nil {
				errMessage := fmt.Sprintf("date must be in the following format:  %s", utils.DateFormat)
				s.RespondErrorDetail(w, r, errMessage, http.StatusBadRequest)
				return
			}

			fetchedSessions, err := s.Queries.ListDoctorSessionsWithinDate(s.Ctx, db.ListDoctorSessionsWithinDateParams{
				DoctorUuid:  doctorUuid,
				StartOfDate: utils.GetStartOfDay(parsedDate),
				EndOfDate:   utils.GetEndOfDay(parsedDate),
			})
			if err != nil || fetchedSessions == nil {
				if err == sql.ErrNoRows {
					s.RespondNoContent(w, r)
					return
				}
				s.RespondErrorStatus(w, r, http.StatusNotFound)
				return
			}

			for _, session := range fetchedSessions {
				sessions = append(sessions, response{
					ID:         session.SessionID,
					GroupIndex: session.SessionGroupIndex,
					Status:     session.SessionStatus,
					Type:       session.SessionType,
					Date:       session.SessionDate.Format(utils.DateFormat),
					Patient: patient{
						Uuid:        session.PatientUuid.String(),
						Name:        session.PatientName,
						PhoneNumber: session.PatientPhoneNumber,
					},
				})
			}
		}

		if len(sessions) < 1 {
			s.RespondNoContent(w, r)
			return
		}

		s.RespondOk(w, r, sessions)
	}
}

func HandleListPatientSessions(s *server.Server) http.HandlerFunc {
	type doctor struct {
		Uuid        string `json:"uuid"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	type patient struct {
		Uuid        string `json:"uuid"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phoneNumber"`
	}
	type response struct {
		ID         int64   `json:"id"`
		Doctor     doctor  `json:"doctor"`
		Patient    patient `json:"patient"`
		GroupIndex int64   `json:"groupIndex"`
		Status     int64   `json:"status"`
		Type       int64   `json:"type"`
		Date       string  `json:"date"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		_, userUuid, err := middlewares.GetAuthDataFromContext(r.Context())
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		patientUuidQueryParam := r.URL.Query().Get("patientUuid")
		patientUuid, err := uuid.FromString(patientUuidQueryParam)
		if err != nil || patientUuid == uuid.Nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		isPatient := userUuid == patientUuid

		upcoming, err := strconv.ParseBool(r.URL.Query().Get("upcoming"))
		if err != nil {
			upcoming = false
		}

		sessions := []response{}

		if upcoming && isPatient {
			fetchedSessions, err := s.Queries.ListUpcomingPatientSessions(s.Ctx, patientUuid)
			if err != nil || fetchedSessions == nil {
				if err == sql.ErrNoRows {
					s.RespondNoContent(w, r)
					return
				}
				s.RespondErrorStatus(w, r, http.StatusNotFound)
				return
			}

			for _, session := range fetchedSessions {
				sessions = append(sessions, response{
					ID:         session.SessionID,
					GroupIndex: session.SessionGroupIndex,
					Status:     session.SessionStatus,
					Type:       session.SessionType,
					Date:       session.SessionDate.Format(utils.DateFormat),
					Doctor: doctor{
						Uuid:        session.DoctorUuid.String(),
						Name:        session.DoctorName,
						Description: session.DoctorDescription,
					},
					Patient: patient{
						Uuid:        session.PatientUuid.String(),
						Name:        session.PatientName,
						PhoneNumber: session.PatientPhoneNumber,
					},
				})
			}
		} else if upcoming && !isPatient {
			fetchedSessions, err := s.Queries.ListUpcomingDoctorPatientSessions(s.Ctx, db.ListUpcomingDoctorPatientSessionsParams{
				DoctorUuid:  userUuid,
				PatientUuid: patientUuid,
			})
			if err != nil || fetchedSessions == nil {
				if err == sql.ErrNoRows {
					s.RespondNoContent(w, r)
					return
				}
				s.RespondErrorStatus(w, r, http.StatusNotFound)
				return
			}

			for _, session := range fetchedSessions {
				sessions = append(sessions, response{
					ID:         session.SessionID,
					GroupIndex: session.SessionGroupIndex,
					Status:     session.SessionStatus,
					Type:       session.SessionType,
					Date:       session.SessionDate.Format(utils.DateFormat),
					Doctor: doctor{
						Uuid:        session.DoctorUuid.String(),
						Name:        session.DoctorName,
						Description: session.DoctorDescription,
					},
					Patient: patient{
						Uuid:        session.PatientUuid.String(),
						Name:        session.PatientName,
						PhoneNumber: session.PatientPhoneNumber,
					},
				})
			}
		} else if !isPatient {
			fetchedSessions, err := s.Queries.ListDoctorPatientSessions(s.Ctx, db.ListDoctorPatientSessionsParams{
				DoctorUuid:  userUuid,
				PatientUuid: patientUuid,
			})
			if err != nil || fetchedSessions == nil {
				if err == sql.ErrNoRows {
					s.RespondNoContent(w, r)
					return
				}
				s.RespondErrorStatus(w, r, http.StatusNotFound)
				return
			}

			for _, session := range fetchedSessions {
				sessions = append(sessions, response{
					ID:         session.SessionID,
					GroupIndex: session.SessionGroupIndex,
					Status:     session.SessionStatus,
					Type:       session.SessionType,
					Date:       session.SessionDate.Format(utils.DateFormat),
					Doctor: doctor{
						Uuid:        session.DoctorUuid.String(),
						Name:        session.DoctorName,
						Description: session.DoctorDescription,
					},
					Patient: patient{
						Uuid:        session.PatientUuid.String(),
						Name:        session.PatientName,
						PhoneNumber: session.PatientPhoneNumber,
					},
				})
			}
		} else {
			fetchedSessions, err := s.Queries.ListPatientSessions(s.Ctx, patientUuid)
			if err != nil || fetchedSessions == nil {
				if err == sql.ErrNoRows {
					s.RespondNoContent(w, r)
					return
				}
				s.RespondErrorStatus(w, r, http.StatusNotFound)
				return
			}

			for _, session := range fetchedSessions {
				sessions = append(sessions, response{
					ID:         session.SessionID,
					GroupIndex: session.SessionGroupIndex,
					Status:     session.SessionStatus,
					Type:       session.SessionType,
					Date:       session.SessionDate.Format(utils.DateFormat),
					Doctor: doctor{
						Uuid:        session.DoctorUuid.String(),
						Name:        session.DoctorName,
						Description: session.DoctorDescription,
					},
					Patient: patient{
						Uuid:        session.PatientUuid.String(),
						Name:        session.PatientName,
						PhoneNumber: session.PatientPhoneNumber,
					},
				})
			}
		}

		if len(sessions) < 1 {
			s.RespondNoContent(w, r)
			return
		}

		s.RespondOk(w, r, sessions)
	}
}
