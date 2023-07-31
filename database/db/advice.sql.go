// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: advice.sql

package db

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

const createAdvice = `-- name: CreateAdvice :one

INSERT INTO advices (message, doctor_uuid) VALUES (?, ?) RETURNING id, message, doctor_uuid, created_at
`

type CreateAdviceParams struct {
	Message    string
	DoctorUuid uuid.UUID
}

func (q *Queries) CreateAdvice(ctx context.Context, arg CreateAdviceParams) (Advice, error) {
	row := q.db.QueryRowContext(ctx, createAdvice, arg.Message, arg.DoctorUuid)
	var i Advice
	err := row.Scan(
		&i.ID,
		&i.Message,
		&i.DoctorUuid,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAdvice = `-- name: DeleteAdvice :exec

DELETE FROM advices WHERE id = ?
`

func (q *Queries) DeleteAdvice(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAdvice, id)
	return err
}

const getAdvice = `-- name: GetAdvice :one

SELECT id, message, doctor_uuid, created_at FROM advices WHERE id = ? LIMIT 1
`

func (q *Queries) GetAdvice(ctx context.Context, id int64) (Advice, error) {
	row := q.db.QueryRowContext(ctx, getAdvice, id)
	var i Advice
	err := row.Scan(
		&i.ID,
		&i.Message,
		&i.DoctorUuid,
		&i.CreatedAt,
	)
	return i, err
}

const listDoctorAdvices = `-- name: ListDoctorAdvices :many

SELECT
    advices.id as advice_id,
    advices.message as advice_message,
    advices.created_at as advice_created_at,
    patients.uuid as patient_uuid,
    doctors.uuid as doctor_uuid,
    doctors.name as doctor_name
FROM advices
    JOIN advice_with_patient ON advices.id = advice_with_patient.advice_id
    JOIN doctors ON advices.doctor_uuid = doctors.uuid
    JOIN patients ON advice_with_patient.patient_uuid = patients.uuid
WHERE advices.doctor_uuid = ?
`

type ListDoctorAdvicesRow struct {
	AdviceID        int64
	AdviceMessage   string
	AdviceCreatedAt time.Time
	PatientUuid     uuid.UUID
	DoctorUuid      uuid.UUID
	DoctorName      string
}

func (q *Queries) ListDoctorAdvices(ctx context.Context, doctorUuid uuid.UUID) ([]ListDoctorAdvicesRow, error) {
	rows, err := q.db.QueryContext(ctx, listDoctorAdvices, doctorUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListDoctorAdvicesRow
	for rows.Next() {
		var i ListDoctorAdvicesRow
		if err := rows.Scan(
			&i.AdviceID,
			&i.AdviceMessage,
			&i.AdviceCreatedAt,
			&i.PatientUuid,
			&i.DoctorUuid,
			&i.DoctorName,
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

const listDoctorPatientAdvices = `-- name: ListDoctorPatientAdvices :many

SELECT
    advices.id as advice_id,
    advices.message as advice_message,
    advices.created_at as advice_created_at,
    patients.uuid as patient_uuid,
    doctors.uuid as doctor_uuid,
    doctors.name as doctor_name
FROM advices
    JOIN advice_with_patient ON advices.id = advice_with_patient.advice_id
    JOIN doctors ON advices.doctor_uuid = doctors.uuid
    JOIN patients ON advice_with_patient.patient_uuid = patients.uuid
WHERE
    advice_with_patient.patient_uuid = ?
    AND advices.doctor_uuid = ?
`

type ListDoctorPatientAdvicesParams struct {
	PatientUuid uuid.UUID
	DoctorUuid  uuid.UUID
}

type ListDoctorPatientAdvicesRow struct {
	AdviceID        int64
	AdviceMessage   string
	AdviceCreatedAt time.Time
	PatientUuid     uuid.UUID
	DoctorUuid      uuid.UUID
	DoctorName      string
}

func (q *Queries) ListDoctorPatientAdvices(ctx context.Context, arg ListDoctorPatientAdvicesParams) ([]ListDoctorPatientAdvicesRow, error) {
	rows, err := q.db.QueryContext(ctx, listDoctorPatientAdvices, arg.PatientUuid, arg.DoctorUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListDoctorPatientAdvicesRow
	for rows.Next() {
		var i ListDoctorPatientAdvicesRow
		if err := rows.Scan(
			&i.AdviceID,
			&i.AdviceMessage,
			&i.AdviceCreatedAt,
			&i.PatientUuid,
			&i.DoctorUuid,
			&i.DoctorName,
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

const listPatientAdvices = `-- name: ListPatientAdvices :many

SELECT
    advices.id as advice_id,
    advices.message as advice_message,
    advices.created_at as advice_created_at,
    patients.uuid as patient_uuid,
    doctors.uuid as doctor_uuid,
    doctors.name as doctor_name
FROM advices
    JOIN advice_with_patient ON advices.id = advice_with_patient.advice_id
    JOIN doctors ON advices.doctor_uuid = doctors.uuid
    JOIN patients ON advice_with_patient.patient_uuid = patients.uuid
WHERE
    advice_with_patient.patient_uuid = ?
GROUP BY advices.id
`

type ListPatientAdvicesRow struct {
	AdviceID        int64
	AdviceMessage   string
	AdviceCreatedAt time.Time
	PatientUuid     uuid.UUID
	DoctorUuid      uuid.UUID
	DoctorName      string
}

func (q *Queries) ListPatientAdvices(ctx context.Context, patientUuid uuid.UUID) ([]ListPatientAdvicesRow, error) {
	rows, err := q.db.QueryContext(ctx, listPatientAdvices, patientUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListPatientAdvicesRow
	for rows.Next() {
		var i ListPatientAdvicesRow
		if err := rows.Scan(
			&i.AdviceID,
			&i.AdviceMessage,
			&i.AdviceCreatedAt,
			&i.PatientUuid,
			&i.DoctorUuid,
			&i.DoctorName,
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
