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
	notConfirmedSessionStatus int64 = iota
	confirmedSessionStatus
	canceledSessionStatus
	concludedSessionStatus
)
const (
	firstSessionStatus = notConfirmedSessionStatus
	_                  = confirmedSessionStatus
	_                  = canceledSessionStatus
	lastSessionStatus  = concludedSessionStatus
)

const (
	notPayedSessionPaymentStatus int64 = iota
	payedSessionPaymentStatus
)

const (
	firstSessionPaymentStatus = notPayedSessionPaymentStatus
	lastSessionPaymentStatus  = payedSessionPaymentStatus
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
		ID            int64   `json:"id"`
		Doctor        doctor  `json:"doctor"`
		Patient       patient `json:"patient"`
		GroupIndex    int64   `json:"groupIndex"`
		Status        int64   `json:"status"`
		PaymentStatus int64   `json:"paymentStatus"`
		Type          int64   `json:"type"`
		Date          string  `json:"date"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, userUuid, err := middlewares.GetAuthDataFromContext(ctx)
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		sessionId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || sessionId < 1 {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		session, err := s.Queries.GetSessionWithParticipants(ctx, int64(sessionId))
		if err != nil || sessionId < 1 {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		if session.DoctorUuid != userUuid && session.PatientUuid != userUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		s.RespondOk(w, r, response{
			ID:            session.SessionID,
			GroupIndex:    session.SessionGroupIndex,
			Status:        session.SessionStatus,
			PaymentStatus: session.SessionPaymentStatus,
			Type:          session.SessionType,
			Date:          session.SessionDate.Format(utils.DateFormat),
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
		DoctorUuid    string `json:"doctorUuid"`
		PatientUuid   string `json:"patientUuid"`
		GroupIndex    int64  `json:"groupIndex"`
		Status        int64  `json:"status"`
		PaymentStatus int64  `json:"paymentStatus"`
		Type          int64  `json:"type"`
		Date          string `json:"date"`
	}
	type response struct {
		ID int64 `json:"id"`
	}
	var sessionAlreadyBookedError = server.Error{
		Type:   "session_already_booked",
		Detail: "There already is a session booked at this time.",
		Status: http.StatusConflict,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, userUuid, err := middlewares.GetAuthDataFromContext(ctx)
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

		if createdSession.Status < firstSessionStatus || createdSession.Status > lastSessionStatus {
			errMessage := fmt.Sprintf("status must be between %d and %d", firstSessionStatus, lastSessionStatus)
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

		_, err = s.Queries.GetDoctorSessionByExactDate(ctx, db.GetDoctorSessionByExactDateParams{
			DoctorUuid: doctorUuid,
			Date:       parsedDate,
		})
		if err == nil {
			s.RespondError(w, r, sessionAlreadyBookedError)
			return
		}

		sessionID, err := s.Queries.CreateSession(ctx, db.CreateSessionParams{
			PatientUuid:   patientUuid,
			DoctorUuid:    doctorUuid,
			Date:          parsedDate,
			GroupIndex:    createdSession.GroupIndex,
			Status:        createdSession.Status,
			PaymentStatus: createdSession.PaymentStatus,
			Type:          createdSession.Type,
		})
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		s.Respond(w, r, response{ID: sessionID}, http.StatusCreated)
	}
}

func HandleUpdateSession(s *server.Server) http.HandlerFunc {
	type request struct {
		Status        *int64  `json:"status"`
		PaymentStatus *int64  `json:"paymentStatus"`
		Date          *string `json:"date"`
	}
	type response struct {
		ID int64 `json:"id"`
	}
	var sessionAlreadyBookedError = server.Error{
		Type:   "session_already_booked",
		Detail: "There already is a session booked at this time.",
		Status: http.StatusConflict,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, userUuid, err := middlewares.GetAuthDataFromContext(ctx)
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

		savedSession, err := s.Queries.GetSession(ctx, int64(sessionId))
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		if savedSession.DoctorUuid != userUuid && savedSession.PatientUuid != userUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		var updateSessionParams = db.UpdateSessionParams{
			ID:            int64(sessionId),
			Status:        sql.NullInt64{Valid: false},
			PaymentStatus: sql.NullInt64{Valid: false},
			Date:          sql.NullTime{Valid: false},
		}

		if updatedSession.Status != nil {
			if *updatedSession.Status < firstSessionPaymentStatus || *updatedSession.PaymentStatus > lastSessionPaymentStatus {
				errMessage := fmt.Sprintf("payment status must be between %d and %d", firstSessionPaymentStatus, lastSessionPaymentStatus)
				s.RespondErrorDetail(w, r, errMessage, http.StatusBadRequest)
				return
			}
			updateSessionParams.PaymentStatus = sql.NullInt64{Int64: *updatedSession.PaymentStatus, Valid: true}
		}

		if updatedSession.PaymentStatus != nil {
			if *updatedSession.PaymentStatus < firstSessionStatus || *updatedSession.Status > lastSessionStatus {
				errMessage := fmt.Sprintf("status must be between %d and %d", firstSessionStatus, lastSessionStatus)
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

			savedSessionWithDate, err := s.Queries.GetDoctorSessionByExactDate(ctx, db.GetDoctorSessionByExactDateParams{
				DoctorUuid: savedSession.DoctorUuid,
				Date:       parsedDate,
			})
			if err == nil && savedSessionWithDate.ID != int64(sessionId) {
				s.RespondError(w, r, sessionAlreadyBookedError)
				return
			}

			updateSessionParams.Date = sql.NullTime{Time: parsedDate, Valid: true}
		}

		updatedSessionID, err := s.Queries.UpdateSession(ctx, updateSessionParams)
		if err != nil {
			if err == sql.ErrNoRows {
				s.RespondErrorStatus(w, r, http.StatusNotFound)
				return
			}
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		s.RespondOk(w, r, response{ID: updatedSessionID})
	}
}

func HandleListSessions(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		doctorUuidQueryParam := r.URL.Query().Get("doctorUuid")
		if strings.TrimSpace(doctorUuidQueryParam) != "" {
			handleListDoctorSessions(s, doctorUuidQueryParam)(w, r)
			return
		}

		patientUuidQueryParam := r.URL.Query().Get("patientUuid")
		if strings.TrimSpace(patientUuidQueryParam) != "" {
			handleListPatientSessions(s, patientUuidQueryParam)(w, r)
			return
		}

		s.RespondErrorStatus(w, r, http.StatusNotFound)
	}
}

func handleListDoctorSessions(s *server.Server, doctorUuidQueryParam string) http.HandlerFunc {
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
		ID            int64   `json:"id"`
		Doctor        doctor  `json:"doctor"`
		Patient       patient `json:"patient"`
		GroupIndex    int64   `json:"groupIndex"`
		Status        int64   `json:"status"`
		PaymentStatus int64   `json:"paymentStatus"`
		Type          int64   `json:"type"`
		Date          string  `json:"date"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, userUuid, err := middlewares.GetAuthDataFromContext(ctx)
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

		dateParam := r.URL.Query().Get("date")

		sessions := []response{}

		fetchedSessions := []db.ListDoctorSessionsRow{}
		if strings.TrimSpace(dateParam) == "" {
			fetchedSessions, err = s.Queries.ListDoctorSessions(ctx, doctorUuid)
		} else {
			var parsedDate time.Time
			parsedDate, err = time.Parse(utils.DateFormat, dateParam)
			if err != nil {
				errMessage := fmt.Sprintf("date must be in the following format:  %s", utils.DateFormat)
				s.RespondErrorDetail(w, r, errMessage, http.StatusBadRequest)
				return
			}

			var list []db.ListDoctorSessionsWithinDateRow
			list, err = s.Queries.ListDoctorSessionsWithinDate(ctx, db.ListDoctorSessionsWithinDateParams{
				DoctorUuid:  doctorUuid,
				StartOfDate: utils.GetStartOfDay(parsedDate),
				EndOfDate:   utils.GetEndOfDay(parsedDate),
			})
			for _, session := range list {
				fetchedSessions = append(fetchedSessions, db.ListDoctorSessionsRow(session))
			}
		}

		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		for _, session := range fetchedSessions {
			sessions = append(sessions, response{
				ID:            session.SessionID,
				GroupIndex:    session.SessionGroupIndex,
				Status:        session.SessionStatus,
				PaymentStatus: session.SessionPaymentStatus,
				Type:          session.SessionType,
				Date:          session.SessionDate.Format(utils.DateFormat),
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

		if len(sessions) < 1 {
			s.RespondNoContent(w, r)
			return
		}

		s.RespondOk(w, r, sessions)
	}
}

func handleListPatientSessions(s *server.Server, patientUuidQueryParam string) http.HandlerFunc {
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
		ID            int64   `json:"id"`
		Doctor        doctor  `json:"doctor"`
		Patient       patient `json:"patient"`
		GroupIndex    int64   `json:"groupIndex"`
		Status        int64   `json:"status"`
		PaymentStatus int64   `json:"paymentStatus"`
		Type          int64   `json:"type"`
		Date          string  `json:"date"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, userUuid, err := middlewares.GetAuthDataFromContext(ctx)
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

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

		fetchedSessions := []db.ListPatientSessionsRow{}
		if upcoming && isPatient {
			var list []db.ListUpcomingPatientSessionsRow
			list, err = s.Queries.ListUpcomingPatientSessions(ctx, patientUuid)
			for _, session := range list {
				fetchedSessions = append(fetchedSessions, db.ListPatientSessionsRow(session))
			}
		} else if upcoming && !isPatient {
			var list []db.ListUpcomingDoctorPatientSessionsRow
			list, err = s.Queries.ListUpcomingDoctorPatientSessions(ctx, db.ListUpcomingDoctorPatientSessionsParams{
				DoctorUuid:  userUuid,
				PatientUuid: patientUuid,
			})
			for _, session := range list {
				fetchedSessions = append(fetchedSessions, db.ListPatientSessionsRow(session))
			}
		} else if !isPatient {
			var list []db.ListDoctorPatientSessionsRow
			list, err = s.Queries.ListDoctorPatientSessions(ctx, db.ListDoctorPatientSessionsParams{
				DoctorUuid:  userUuid,
				PatientUuid: patientUuid,
			})
			for _, session := range list {
				fetchedSessions = append(fetchedSessions, db.ListPatientSessionsRow(session))
			}
		} else {
			fetchedSessions, err = s.Queries.ListPatientSessions(ctx, patientUuid)
		}

		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		for _, session := range fetchedSessions {
			sessions = append(sessions, response{
				ID:            session.SessionID,
				GroupIndex:    session.SessionGroupIndex,
				Status:        session.SessionStatus,
				PaymentStatus: session.SessionPaymentStatus,
				Type:          session.SessionType,
				Date:          session.SessionDate.Format(utils.DateFormat),
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

		if len(sessions) < 1 {
			s.RespondNoContent(w, r)
			return
		}

		s.RespondOk(w, r, sessions)
	}
}
