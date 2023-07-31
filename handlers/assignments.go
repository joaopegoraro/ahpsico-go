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

const (
	pending int64 = iota
	done
	missed
)
const (
	firstAssignmentStatus = pending
	_                     = done
	lastAssignmentStatus  = missed
)

func HandleCreateAssignment(s *server.Server) http.HandlerFunc {
	type request struct {
		DoctorUuid        string `json:"doctorUuid"`
		PatientUuid       string `json:"patientUuid"`
		DeliverySessionID int64  `json:"deliverySessionId"`
		Title             string `json:"title"`
		Description       string `json:"description"`
		Status            int64  `json:"status"`
	}
	type response struct {
		ID                int64  `json:"id"`
		DoctorUuid        string `json:"doctorUuid"`
		PatientUuid       string `json:"patientUuid"`
		DeliverySessionID int64  `json:"deliverySessionId"`
		Title             string `json:"title"`
		Description       string `json:"description"`
		Status            int64  `json:"status"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		_, userUuid, err := middlewares.GetAuthDataFromContext(r.Context())
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		var createdAssignment request
		err = s.Decode(w, r, &createdAssignment)
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		doctorUuid, err := uuid.FromString(createdAssignment.DoctorUuid)
		if err != nil {
			s.RespondErrorDetail(w, r, "invalid doctor uuid", http.StatusBadRequest)
			return
		}

		patientUuid, err := uuid.FromString(createdAssignment.PatientUuid)
		if err != nil {
			s.RespondErrorDetail(w, r, "invalid patient uuid", http.StatusBadRequest)
			return
		}

		if createdAssignment.DeliverySessionID < 1 {
			s.RespondErrorDetail(w, r, "invalid session id", http.StatusBadRequest)
			return
		}

		if doctorUuid != userUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		if strings.TrimSpace(createdAssignment.Title) == "" {
			s.RespondErrorDetail(w, r, "title cannot be empty", http.StatusBadRequest)
			return
		}

		if createdAssignment.Status < firstAssignmentStatus || createdAssignment.Status > lastAssignmentStatus {
			errMessage := fmt.Sprintf("status must be between %d and %d", firstAssignmentStatus, lastAssignmentStatus)
			s.RespondErrorDetail(w, r, errMessage, http.StatusBadRequest)
			return
		}

		assignment, err := s.Queries.CreateAssignment(s.Ctx, db.CreateAssignmentParams{
			Title:       createdAssignment.Title,
			Description: createdAssignment.Description,
			PatientUuid: patientUuid,
			DoctorUuid:  doctorUuid,
			SessionID:   createdAssignment.DeliverySessionID,
			Status:      createdAssignment.Status,
		})
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		s.Respond(w, r, response{
			ID:                assignment.ID,
			Title:             assignment.Title,
			Description:       assignment.Description,
			DoctorUuid:        assignment.DoctorUuid.String(),
			DeliverySessionID: assignment.SessionID,
			PatientUuid:       assignment.PatientUuid.String(),
			Status:            assignment.Status,
		}, http.StatusCreated)
	}
}

func HandleDeleteAssignment(s *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, userUuid, err := middlewares.GetAuthDataFromContext(r.Context())
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		assignmentID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || assignmentID < 1 {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		savedAssignment, err := s.Queries.GetAssignment(s.Ctx, int64(assignmentID))
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		if savedAssignment.DoctorUuid != userUuid && savedAssignment.PatientUuid != userUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		err = s.Queries.DeleteAssignment(s.Ctx, int64(assignmentID))
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusInternalServerError)
			return
		}

		s.RespondNoContent(w, r)
	}
}

func HandleUpdateAssignment(s *server.Server) http.HandlerFunc {
	type request struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Status      *int64  `json:"status"`
	}
	type response struct {
		ID                int64  `json:"id"`
		DoctorUuid        string `json:"doctorUuid"`
		PatientUuid       string `json:"patientUuid"`
		DeliverySessionID int64  `json:"deliverySessionId"`
		Title             string `json:"title"`
		Description       string `json:"description"`
		Status            int64  `json:"status"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		_, userUuid, err := middlewares.GetAuthDataFromContext(r.Context())
		if err != nil {
			middlewares.RespondAuthError(w, r, s)
			return
		}

		assignmentID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || assignmentID < 1 {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		var updatedAssignment request
		err = s.Decode(w, r, &updatedAssignment)
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		savedAssignment, err := s.Queries.GetSession(s.Ctx, int64(assignmentID))
		if err != nil {
			s.RespondErrorStatus(w, r, http.StatusNotFound)
			return
		}

		if savedAssignment.DoctorUuid != userUuid && savedAssignment.PatientUuid != userUuid {
			s.RespondErrorStatus(w, r, http.StatusForbidden)
			return
		}

		var updateAssignmentParams = db.UpdateAssignmentParams{
			ID:          int64(assignmentID),
			Title:       sql.NullString{Valid: false},
			Description: sql.NullString{Valid: false},
			Status:      sql.NullInt64{Valid: false},
		}

		if updatedAssignment.Title != nil {
			if strings.TrimSpace(*updatedAssignment.Title) == "" {
				s.RespondErrorDetail(w, r, "title cannot be blank", http.StatusBadRequest)
				return
			}
			updateAssignmentParams.Title = sql.NullString{String: *updatedAssignment.Title, Valid: true}
		}

		if updatedAssignment.Description != nil {
			if strings.TrimSpace(*updatedAssignment.Description) == "" {
				s.RespondErrorDetail(w, r, "description cannot be blank", http.StatusBadRequest)
				return
			}
			updateAssignmentParams.Description = sql.NullString{String: *updatedAssignment.Description, Valid: true}
		}

		if updatedAssignment.Status != nil {
			if *updatedAssignment.Status < firstAssignmentStatus || *updatedAssignment.Status > lastAssignmentStatus {
				errMessage := fmt.Sprintf("status must be between %d and %d", firstAssignmentStatus, lastAssignmentStatus)
				s.RespondErrorDetail(w, r, errMessage, http.StatusBadRequest)
				return
			}
			updateAssignmentParams.Status = sql.NullInt64{Int64: *updatedAssignment.Status, Valid: true}
		}

		assignment, err := s.Queries.UpdateAssignment(s.Ctx, updateAssignmentParams)
		if err != nil {
			if err == sql.ErrNoRows {
				s.RespondErrorStatus(w, r, http.StatusNotFound)
				return
			}
			s.RespondErrorStatus(w, r, http.StatusBadRequest)
			return
		}

		s.RespondOk(w, r, response{
			ID:                assignment.ID,
			Title:             assignment.Title,
			Description:       assignment.Description,
			DoctorUuid:        assignment.DoctorUuid.String(),
			DeliverySessionID: assignment.SessionID,
			PatientUuid:       assignment.PatientUuid.String(),
			Status:            assignment.Status,
		})
	}
}

func HandleListPatientAssignments(s *server.Server) http.HandlerFunc {
	type doctor struct {
		Uuid        string `json:"uuid"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	type session struct {
		ID   int64  `json:"id"`
		Date string `json:"date"`
	}
	type response struct {
		ID              int64   `json:"id"`
		Doctor          doctor  `json:"doctor"`
		PatientUuid     string  `json:"patientUuid"`
		DeliverySession session `json:"deliverySession"`
		Title           string  `json:"title"`
		Description     string  `json:"description"`
		Status          int64   `json:"status"`
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

		pending, err := strconv.ParseBool(r.URL.Query().Get("pending"))
		if err != nil {
			pending = false
		}

		assignments := []response{}

		if pending && isPatient {
			fetchedAssignments, err := s.Queries.ListPendingPatientAssignments(s.Ctx, patientUuid)
			if err != nil || fetchedAssignments == nil {
				if err == sql.ErrNoRows {
					s.RespondNoContent(w, r)
					return
				}
				s.RespondErrorStatus(w, r, http.StatusNotFound)
				return
			}

			for _, assignment := range fetchedAssignments {
				assignments = append(assignments, response{
					ID:          assignment.AssignmentID,
					Title:       assignment.AssignmentTitle,
					Description: assignment.AssignmentDescription,
					Status:      assignment.AssignmentStatus,
					Doctor: doctor{
						Uuid:        assignment.DoctorUuid.String(),
						Name:        assignment.DoctorName,
						Description: assignment.DoctorDescription,
					},
					PatientUuid: assignment.PatientUuid.String(),
					DeliverySession: session{
						ID:   assignment.SessionID,
						Date: assignment.SessionDate.Format(utils.DateFormat),
					},
				})
			}
		} else if pending && !isPatient {
			fetchedAssignments, err := s.Queries.ListPendingDoctorPatientAssignments(s.Ctx, db.ListPendingDoctorPatientAssignmentsParams{
				DoctorUuid:  userUuid,
				PatientUuid: patientUuid,
			})
			if err != nil || fetchedAssignments == nil {
				if err == sql.ErrNoRows {
					s.RespondNoContent(w, r)
					return
				}
				s.RespondErrorStatus(w, r, http.StatusNotFound)
				return
			}

			for _, assignment := range fetchedAssignments {
				assignments = append(assignments, response{
					ID:          assignment.AssignmentID,
					Title:       assignment.AssignmentTitle,
					Description: assignment.AssignmentDescription,
					Status:      assignment.AssignmentStatus,
					Doctor: doctor{
						Uuid:        assignment.DoctorUuid.String(),
						Name:        assignment.DoctorName,
						Description: assignment.DoctorDescription,
					},
					PatientUuid: assignment.PatientUuid.String(),
					DeliverySession: session{
						ID:   assignment.SessionID,
						Date: assignment.SessionDate.Format(utils.DateFormat),
					},
				})
			}
		} else if !isPatient {
			fetchedAssignments, err := s.Queries.ListDoctorPatientAssignments(s.Ctx, db.ListDoctorPatientAssignmentsParams{
				DoctorUuid:  userUuid,
				PatientUuid: patientUuid,
			})
			if err != nil || fetchedAssignments == nil {
				if err == sql.ErrNoRows {
					s.RespondNoContent(w, r)
					return
				}
				s.RespondErrorStatus(w, r, http.StatusNotFound)
				return
			}

			for _, assignment := range fetchedAssignments {
				assignments = append(assignments, response{
					ID:          assignment.AssignmentID,
					Title:       assignment.AssignmentTitle,
					Description: assignment.AssignmentDescription,
					Status:      assignment.AssignmentStatus,
					Doctor: doctor{
						Uuid:        assignment.DoctorUuid.String(),
						Name:        assignment.DoctorName,
						Description: assignment.DoctorDescription,
					},
					PatientUuid: assignment.PatientUuid.String(),
					DeliverySession: session{
						ID:   assignment.SessionID,
						Date: assignment.SessionDate.Format(utils.DateFormat),
					},
				})
			}
		} else {
			fetchedAssignments, err := s.Queries.ListPatientAssignments(s.Ctx, patientUuid)
			if err != nil || fetchedAssignments == nil {
				if err == sql.ErrNoRows {
					s.RespondNoContent(w, r)
					return
				}
				s.RespondErrorStatus(w, r, http.StatusNotFound)
				return
			}

			for _, assignment := range fetchedAssignments {
				assignments = append(assignments, response{
					ID:          assignment.AssignmentID,
					Title:       assignment.AssignmentTitle,
					Description: assignment.AssignmentDescription,
					Status:      assignment.AssignmentStatus,
					Doctor: doctor{
						Uuid:        assignment.DoctorUuid.String(),
						Name:        assignment.DoctorName,
						Description: assignment.DoctorDescription,
					},
					PatientUuid: assignment.PatientUuid.String(),
					DeliverySession: session{
						ID:   assignment.SessionID,
						Date: assignment.SessionDate.Format(utils.DateFormat),
					},
				})
			}
		}

		if len(assignments) < 1 {
			s.RespondNoContent(w, r)
			return
		}

		s.RespondOk(w, r, assignments)
	}
}
