// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: session.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

const createSession = `-- name: CreateSession :one

INSERT INTO
    sessions (
        patient_uuid,
        doctor_uuid,
        date,
        group_index,
        type,
        status
    )
VALUES (?, ?, ?, ?, ?, ?) RETURNING id, patient_uuid, doctor_uuid, date, group_index, type, status, created_at, updated_at
`

type CreateSessionParams struct {
	PatientUuid uuid.UUID
	DoctorUuid  uuid.UUID
	Date        time.Time
	GroupIndex  int64
	Type        int64
	Status      int64
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, createSession,
		arg.PatientUuid,
		arg.DoctorUuid,
		arg.Date,
		arg.GroupIndex,
		arg.Type,
		arg.Status,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.PatientUuid,
		&i.DoctorUuid,
		&i.Date,
		&i.GroupIndex,
		&i.Type,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteSession = `-- name: DeleteSession :exec

DELETE FROM sessions where id = ?
`

func (q *Queries) DeleteSession(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteSession, id)
	return err
}

const getDoctorSessionByExactDate = `-- name: GetDoctorSessionByExactDate :one

SELECT sessions.id, sessions.patient_uuid, sessions.doctor_uuid, sessions.date, sessions.group_index, sessions.type, sessions.status, sessions.created_at, sessions.updated_at
FROM sessions
WHERE doctor_uuid = ? AND date = ?
LIMIT 1
`

type GetDoctorSessionByExactDateParams struct {
	DoctorUuid uuid.UUID
	Date       time.Time
}

func (q *Queries) GetDoctorSessionByExactDate(ctx context.Context, arg GetDoctorSessionByExactDateParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, getDoctorSessionByExactDate, arg.DoctorUuid, arg.Date)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.PatientUuid,
		&i.DoctorUuid,
		&i.Date,
		&i.GroupIndex,
		&i.Type,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSession = `-- name: GetSession :one

SELECT id, patient_uuid, doctor_uuid, date, group_index, type, status, created_at, updated_at FROM sessions WHERE id = ? LIMIT 1
`

func (q *Queries) GetSession(ctx context.Context, id int64) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSession, id)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.PatientUuid,
		&i.DoctorUuid,
		&i.Date,
		&i.GroupIndex,
		&i.Type,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSessionWithParticipants = `-- name: GetSessionWithParticipants :one

SELECT
    s.id as session_id,
    s.date as session_date,
    s.group_index as session_group_index,
    s.type as session_type,
    s.status as session_status,
    s.created_at as session_created_at,
    d.uuid as doctor_uuid,
    d.name as doctor_name,
    d.description as doctor_description,
    p.uuid as patient_uuid,
    p.name as patient_name,
    p.phone_number as patient_phone_number
FROM sessions s
    JOIN doctors d ON doctors.uuid = sessions.doctor_uuid
    JOIN patients p ON patients.uuid = sessions.patient_uuid
WHERE id = ?
LIMIT 1
`

type GetSessionWithParticipantsRow struct {
	SessionID          int64
	SessionDate        time.Time
	SessionGroupIndex  int64
	SessionType        int64
	SessionStatus      int64
	SessionCreatedAt   time.Time
	DoctorUuid         uuid.UUID
	DoctorName         string
	DoctorDescription  string
	PatientUuid        uuid.UUID
	PatientName        string
	PatientPhoneNumber string
}

func (q *Queries) GetSessionWithParticipants(ctx context.Context, id int64) (GetSessionWithParticipantsRow, error) {
	row := q.db.QueryRowContext(ctx, getSessionWithParticipants, id)
	var i GetSessionWithParticipantsRow
	err := row.Scan(
		&i.SessionID,
		&i.SessionDate,
		&i.SessionGroupIndex,
		&i.SessionType,
		&i.SessionStatus,
		&i.SessionCreatedAt,
		&i.DoctorUuid,
		&i.DoctorName,
		&i.DoctorDescription,
		&i.PatientUuid,
		&i.PatientName,
		&i.PatientPhoneNumber,
	)
	return i, err
}

const listDoctorPatientSessions = `-- name: ListDoctorPatientSessions :many

SELECT
    s.id as session_id,
    s.date as session_date,
    s.group_index as session_group_index,
    s.type as session_type,
    s.status as session_status,
    s.created_at as session_created_at,
    d.uuid as doctor_uuid,
    d.name as doctor_name,
    d.description as doctor_description,
    p.uuid as patient_uuid,
    p.name as patient_name,
    p.phone_number as patient_phone_number
FROM sessions s
    JOIN doctors d ON doctors.uuid = sessions.doctor_uuid
    JOIN patients p ON patients.uuid = sessions.patient_uuid
WHERE
    patient_uuid = ?
    AND doctor_uuid = ?
`

type ListDoctorPatientSessionsParams struct {
	PatientUuid uuid.UUID
	DoctorUuid  uuid.UUID
}

type ListDoctorPatientSessionsRow struct {
	SessionID          int64
	SessionDate        time.Time
	SessionGroupIndex  int64
	SessionType        int64
	SessionStatus      int64
	SessionCreatedAt   time.Time
	DoctorUuid         uuid.UUID
	DoctorName         string
	DoctorDescription  string
	PatientUuid        uuid.UUID
	PatientName        string
	PatientPhoneNumber string
}

func (q *Queries) ListDoctorPatientSessions(ctx context.Context, arg ListDoctorPatientSessionsParams) ([]ListDoctorPatientSessionsRow, error) {
	rows, err := q.db.QueryContext(ctx, listDoctorPatientSessions, arg.PatientUuid, arg.DoctorUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListDoctorPatientSessionsRow
	for rows.Next() {
		var i ListDoctorPatientSessionsRow
		if err := rows.Scan(
			&i.SessionID,
			&i.SessionDate,
			&i.SessionGroupIndex,
			&i.SessionType,
			&i.SessionStatus,
			&i.SessionCreatedAt,
			&i.DoctorUuid,
			&i.DoctorName,
			&i.DoctorDescription,
			&i.PatientUuid,
			&i.PatientName,
			&i.PatientPhoneNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listDoctorSessions = `-- name: ListDoctorSessions :many

SELECT
    s.id as session_id,
    s.date as session_date,
    s.group_index as session_group_index,
    s.type as session_type,
    s.status as session_status,
    s.created_at as session_created_at,
    p.uuid as patient_uuid,
    p.name as patient_name,
    p.phone_number as patient_phone_number
FROM sessions s
    JOIN patients p ON patients.uuid = sessions.patient_uuid
WHERE doctor_uuid = ?
`

type ListDoctorSessionsRow struct {
	SessionID          int64
	SessionDate        time.Time
	SessionGroupIndex  int64
	SessionType        int64
	SessionStatus      int64
	SessionCreatedAt   time.Time
	PatientUuid        uuid.UUID
	PatientName        string
	PatientPhoneNumber string
}

func (q *Queries) ListDoctorSessions(ctx context.Context, doctorUuid uuid.UUID) ([]ListDoctorSessionsRow, error) {
	rows, err := q.db.QueryContext(ctx, listDoctorSessions, doctorUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListDoctorSessionsRow
	for rows.Next() {
		var i ListDoctorSessionsRow
		if err := rows.Scan(
			&i.SessionID,
			&i.SessionDate,
			&i.SessionGroupIndex,
			&i.SessionType,
			&i.SessionStatus,
			&i.SessionCreatedAt,
			&i.PatientUuid,
			&i.PatientName,
			&i.PatientPhoneNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listDoctorSessionsWithinDate = `-- name: ListDoctorSessionsWithinDate :many

SELECT
    s.id as session_id,
    s.date as session_date,
    s.group_index as session_group_index,
    s.type as session_type,
    s.status as session_status,
    s.created_at as session_created_at,
    d.uuid as doctor_uuid,
    d.name as doctor_name,
    d.description as doctor_description,
    p.uuid as patient_uuid,
    p.name as patient_name,
    p.phone_number as patient_phone_number
FROM sessions s
    JOIN doctors d ON doctors.uuid = sessions.doctor_uuid
    JOIN patients p ON patients.uuid = sessions.patient_uuid
WHERE
    doctor_uuid = ?1
    AND date >= ?2
    AND date <= ?3
`

type ListDoctorSessionsWithinDateParams struct {
	DoctorUuid  uuid.UUID
	StartOfDate time.Time
	EndOfDate   time.Time
}

type ListDoctorSessionsWithinDateRow struct {
	SessionID          int64
	SessionDate        time.Time
	SessionGroupIndex  int64
	SessionType        int64
	SessionStatus      int64
	SessionCreatedAt   time.Time
	DoctorUuid         uuid.UUID
	DoctorName         string
	DoctorDescription  string
	PatientUuid        uuid.UUID
	PatientName        string
	PatientPhoneNumber string
}

func (q *Queries) ListDoctorSessionsWithinDate(ctx context.Context, arg ListDoctorSessionsWithinDateParams) ([]ListDoctorSessionsWithinDateRow, error) {
	rows, err := q.db.QueryContext(ctx, listDoctorSessionsWithinDate, arg.DoctorUuid, arg.StartOfDate, arg.EndOfDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListDoctorSessionsWithinDateRow
	for rows.Next() {
		var i ListDoctorSessionsWithinDateRow
		if err := rows.Scan(
			&i.SessionID,
			&i.SessionDate,
			&i.SessionGroupIndex,
			&i.SessionType,
			&i.SessionStatus,
			&i.SessionCreatedAt,
			&i.DoctorUuid,
			&i.DoctorName,
			&i.DoctorDescription,
			&i.PatientUuid,
			&i.PatientName,
			&i.PatientPhoneNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPatientSessions = `-- name: ListPatientSessions :many

SELECT
    s.id as session_id,
    s.date as session_date,
    s.group_index as session_group_index,
    s.type as session_type,
    s.status as session_status,
    s.created_at as session_created_at,
    d.uuid as doctor_uuid,
    d.name as doctor_name,
    d.description as doctor_description,
    p.uuid as patient_uuid,
    p.name as patient_name,
    p.phone_number as patient_phone_number
FROM sessions s
    JOIN doctors d ON doctors.uuid = sessions.doctor_uuid
    JOIN patients p ON patients.uuid = sessions.patient_uuid
WHERE patient_uuid = ?
`

type ListPatientSessionsRow struct {
	SessionID          int64
	SessionDate        time.Time
	SessionGroupIndex  int64
	SessionType        int64
	SessionStatus      int64
	SessionCreatedAt   time.Time
	DoctorUuid         uuid.UUID
	DoctorName         string
	DoctorDescription  string
	PatientUuid        uuid.UUID
	PatientName        string
	PatientPhoneNumber string
}

func (q *Queries) ListPatientSessions(ctx context.Context, patientUuid uuid.UUID) ([]ListPatientSessionsRow, error) {
	rows, err := q.db.QueryContext(ctx, listPatientSessions, patientUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListPatientSessionsRow
	for rows.Next() {
		var i ListPatientSessionsRow
		if err := rows.Scan(
			&i.SessionID,
			&i.SessionDate,
			&i.SessionGroupIndex,
			&i.SessionType,
			&i.SessionStatus,
			&i.SessionCreatedAt,
			&i.DoctorUuid,
			&i.DoctorName,
			&i.DoctorDescription,
			&i.PatientUuid,
			&i.PatientName,
			&i.PatientPhoneNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUpcomingDoctorPatientSessions = `-- name: ListUpcomingDoctorPatientSessions :many

SELECT
    s.id as session_id,
    s.date as session_date,
    s.group_index as session_group_index,
    s.type as session_type,
    s.status as session_status,
    s.created_at as session_created_at,
    d.uuid as doctor_uuid,
    d.name as doctor_name,
    d.description as doctor_description,
    p.uuid as patient_uuid,
    p.name as patient_name,
    p.phone_number as patient_phone_number
FROM sessions s
    JOIN doctors d ON doctors.uuid = sessions.doctor_uuid
    JOIN patients p ON patients.uuid = sessions.patient_uuid
WHERE
    patient_uuid = ?
    AND doctor_uuid = ?
    AND date >= CURRENT_TIMESTAMP
`

type ListUpcomingDoctorPatientSessionsParams struct {
	PatientUuid uuid.UUID
	DoctorUuid  uuid.UUID
}

type ListUpcomingDoctorPatientSessionsRow struct {
	SessionID          int64
	SessionDate        time.Time
	SessionGroupIndex  int64
	SessionType        int64
	SessionStatus      int64
	SessionCreatedAt   time.Time
	DoctorUuid         uuid.UUID
	DoctorName         string
	DoctorDescription  string
	PatientUuid        uuid.UUID
	PatientName        string
	PatientPhoneNumber string
}

func (q *Queries) ListUpcomingDoctorPatientSessions(ctx context.Context, arg ListUpcomingDoctorPatientSessionsParams) ([]ListUpcomingDoctorPatientSessionsRow, error) {
	rows, err := q.db.QueryContext(ctx, listUpcomingDoctorPatientSessions, arg.PatientUuid, arg.DoctorUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUpcomingDoctorPatientSessionsRow
	for rows.Next() {
		var i ListUpcomingDoctorPatientSessionsRow
		if err := rows.Scan(
			&i.SessionID,
			&i.SessionDate,
			&i.SessionGroupIndex,
			&i.SessionType,
			&i.SessionStatus,
			&i.SessionCreatedAt,
			&i.DoctorUuid,
			&i.DoctorName,
			&i.DoctorDescription,
			&i.PatientUuid,
			&i.PatientName,
			&i.PatientPhoneNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUpcomingPatientSessions = `-- name: ListUpcomingPatientSessions :many

SELECT
    s.id as session_id,
    s.date as session_date,
    s.group_index as session_group_index,
    s.type as session_type,
    s.status as session_status,
    s.created_at as session_created_at,
    d.uuid as doctor_uuid,
    d.name as doctor_name,
    d.description as doctor_description,
    p.uuid as patient_uuid,
    p.name as patient_name,
    p.phone_number as patient_phone_number
FROM sessions s
    JOIN doctors d ON doctors.uuid = sessions.doctor_uuid
    JOIN patients p ON patients.uuid = sessions.patient_uuid
WHERE
    patient_uuid = ?
    AND date >= CURRENT_TIMESTAMP
`

type ListUpcomingPatientSessionsRow struct {
	SessionID          int64
	SessionDate        time.Time
	SessionGroupIndex  int64
	SessionType        int64
	SessionStatus      int64
	SessionCreatedAt   time.Time
	DoctorUuid         uuid.UUID
	DoctorName         string
	DoctorDescription  string
	PatientUuid        uuid.UUID
	PatientName        string
	PatientPhoneNumber string
}

func (q *Queries) ListUpcomingPatientSessions(ctx context.Context, patientUuid uuid.UUID) ([]ListUpcomingPatientSessionsRow, error) {
	rows, err := q.db.QueryContext(ctx, listUpcomingPatientSessions, patientUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUpcomingPatientSessionsRow
	for rows.Next() {
		var i ListUpcomingPatientSessionsRow
		if err := rows.Scan(
			&i.SessionID,
			&i.SessionDate,
			&i.SessionGroupIndex,
			&i.SessionType,
			&i.SessionStatus,
			&i.SessionCreatedAt,
			&i.DoctorUuid,
			&i.DoctorName,
			&i.DoctorDescription,
			&i.PatientUuid,
			&i.PatientName,
			&i.PatientPhoneNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateSession = `-- name: UpdateSession :one

UPDATE sessions
SET
    date = COALESCE(?1, date),
    status = COALESCE(?2, status),
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = ?3 RETURNING id, patient_uuid, doctor_uuid, date, group_index, type, status, created_at, updated_at
`

type UpdateSessionParams struct {
	Date   sql.NullTime
	Status sql.NullInt64
	ID     int64
}

func (q *Queries) UpdateSession(ctx context.Context, arg UpdateSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, updateSession, arg.Date, arg.Status, arg.ID)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.PatientUuid,
		&i.DoctorUuid,
		&i.Date,
		&i.GroupIndex,
		&i.Type,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
