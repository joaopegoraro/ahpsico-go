// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: assignment.sql

package db

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
)

const createAssignment = `-- name: CreateAssignment :one

INSERT INTO
    assignments (
        title,
        description,
        patient_uuid,
        doctor_uuid,
        session_id,
        status
    )
VALUES (?, ?, ?, ?, ?, ?) RETURNING id, title, description, patient_uuid, doctor_uuid, session_id, status, created_at, updated_at
`

type CreateAssignmentParams struct {
	Title       string
	Description string
	PatientUuid uuid.UUID
	DoctorUuid  uuid.UUID
	SessionID   int64
	Status      int64
}

func (q *Queries) CreateAssignment(ctx context.Context, arg CreateAssignmentParams) (Assignment, error) {
	row := q.db.QueryRowContext(ctx, createAssignment,
		arg.Title,
		arg.Description,
		arg.PatientUuid,
		arg.DoctorUuid,
		arg.SessionID,
		arg.Status,
	)
	var i Assignment
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.PatientUuid,
		&i.DoctorUuid,
		&i.SessionID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAssignment = `-- name: DeleteAssignment :exec

DELETE FROM assignments where id = ?
`

func (q *Queries) DeleteAssignment(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAssignment, id)
	return err
}

const listDoctorPatientAssignments = `-- name: ListDoctorPatientAssignments :many

SELECT assignments.id, assignments.title, assignments.description, assignments.patient_uuid, assignments.doctor_uuid, assignments.session_id, assignments.status, assignments.created_at, assignments.updated_at
FROM assignments
WHERE
    patient_uuid = ?
    AND doctor_uuid = ?
`

type ListDoctorPatientAssignmentsParams struct {
	PatientUuid uuid.UUID
	DoctorUuid  uuid.UUID
}

func (q *Queries) ListDoctorPatientAssignments(ctx context.Context, arg ListDoctorPatientAssignmentsParams) ([]Assignment, error) {
	rows, err := q.db.QueryContext(ctx, listDoctorPatientAssignments, arg.PatientUuid, arg.DoctorUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Assignment
	for rows.Next() {
		var i Assignment
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.PatientUuid,
			&i.DoctorUuid,
			&i.SessionID,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const listPatientAssignments = `-- name: ListPatientAssignments :many

SELECT assignments.id, assignments.title, assignments.description, assignments.patient_uuid, assignments.doctor_uuid, assignments.session_id, assignments.status, assignments.created_at, assignments.updated_at FROM assignments WHERE patient_uuid = ?
`

func (q *Queries) ListPatientAssignments(ctx context.Context, patientUuid uuid.UUID) ([]Assignment, error) {
	rows, err := q.db.QueryContext(ctx, listPatientAssignments, patientUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Assignment
	for rows.Next() {
		var i Assignment
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.PatientUuid,
			&i.DoctorUuid,
			&i.SessionID,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const listPendingDoctorPatientAssignments = `-- name: ListPendingDoctorPatientAssignments :many

SELECT assignments.id, assignments.title, assignments.description, assignments.patient_uuid, assignments.doctor_uuid, assignments.session_id, assignments.status, assignments.created_at, assignments.updated_at
FROM assignments
WHERE
    patient_uuid = ?
    AND doctor_uuid = ?
    AND status = 0
`

type ListPendingDoctorPatientAssignmentsParams struct {
	PatientUuid uuid.UUID
	DoctorUuid  uuid.UUID
}

func (q *Queries) ListPendingDoctorPatientAssignments(ctx context.Context, arg ListPendingDoctorPatientAssignmentsParams) ([]Assignment, error) {
	rows, err := q.db.QueryContext(ctx, listPendingDoctorPatientAssignments, arg.PatientUuid, arg.DoctorUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Assignment
	for rows.Next() {
		var i Assignment
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.PatientUuid,
			&i.DoctorUuid,
			&i.SessionID,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const listPendingPatientAssignments = `-- name: ListPendingPatientAssignments :many

SELECT assignments.id, assignments.title, assignments.description, assignments.patient_uuid, assignments.doctor_uuid, assignments.session_id, assignments.status, assignments.created_at, assignments.updated_at
FROM assignments
WHERE
    patient_uuid = ?
    AND status = 0
`

func (q *Queries) ListPendingPatientAssignments(ctx context.Context, patientUuid uuid.UUID) ([]Assignment, error) {
	rows, err := q.db.QueryContext(ctx, listPendingPatientAssignments, patientUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Assignment
	for rows.Next() {
		var i Assignment
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.PatientUuid,
			&i.DoctorUuid,
			&i.SessionID,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updateAssignment = `-- name: UpdateAssignment :one

UPDATE assignments
SET
    title = COALESCE(?1, title),
    description = COALESCE(
        ?2,
        description
    ),
    status = COALESCE(?3, status),
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = ?4 RETURNING id, title, description, patient_uuid, doctor_uuid, session_id, status, created_at, updated_at
`

type UpdateAssignmentParams struct {
	Title       sql.NullString
	Description sql.NullString
	Status      sql.NullInt64
	ID          int64
}

func (q *Queries) UpdateAssignment(ctx context.Context, arg UpdateAssignmentParams) (Assignment, error) {
	row := q.db.QueryRowContext(ctx, updateAssignment,
		arg.Title,
		arg.Description,
		arg.Status,
		arg.ID,
	)
	var i Assignment
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.PatientUuid,
		&i.DoctorUuid,
		&i.SessionID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
