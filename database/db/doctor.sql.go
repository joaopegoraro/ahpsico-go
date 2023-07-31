// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: doctor.sql

package db

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
)

const createDoctor = `-- name: CreateDoctor :one

INSERT INTO
    doctors (
        uuid,
        name,
        phone_number,
        description,
        crp,
        pix_key,
        payment_details
    )
VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING uuid, name, phone_number, description, crp, pix_key, payment_details, created_at, updated_at
`

type CreateDoctorParams struct {
	Uuid           uuid.UUID
	Name           string
	PhoneNumber    string
	Description    string
	Crp            string
	PixKey         string
	PaymentDetails string
}

func (q *Queries) CreateDoctor(ctx context.Context, arg CreateDoctorParams) (Doctor, error) {
	row := q.db.QueryRowContext(ctx, createDoctor,
		arg.Uuid,
		arg.Name,
		arg.PhoneNumber,
		arg.Description,
		arg.Crp,
		arg.PixKey,
		arg.PaymentDetails,
	)
	var i Doctor
	err := row.Scan(
		&i.Uuid,
		&i.Name,
		&i.PhoneNumber,
		&i.Description,
		&i.Crp,
		&i.PixKey,
		&i.PaymentDetails,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getDoctor = `-- name: GetDoctor :one

SELECT uuid, name, phone_number, description, crp, pix_key, payment_details, created_at, updated_at FROM doctors WHERE uuid = ? LIMIT 1
`

func (q *Queries) GetDoctor(ctx context.Context, argUuid uuid.UUID) (Doctor, error) {
	row := q.db.QueryRowContext(ctx, getDoctor, argUuid)
	var i Doctor
	err := row.Scan(
		&i.Uuid,
		&i.Name,
		&i.PhoneNumber,
		&i.Description,
		&i.Crp,
		&i.PixKey,
		&i.PaymentDetails,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listPatientDoctors = `-- name: ListPatientDoctors :many

SELECT doctors.uuid, doctors.name, doctors.phone_number, doctors.description, doctors.crp, doctors.pix_key, doctors.payment_details, doctors.created_at, doctors.updated_at
FROM doctors
    JOIN patient_with_doctor ON doctors.uuid = patient_with_doctor.doctor_uuid
WHERE
    patient_with_doctor.patient_uuid = ?
`

func (q *Queries) ListPatientDoctors(ctx context.Context, patientUuid uuid.UUID) ([]Doctor, error) {
	rows, err := q.db.QueryContext(ctx, listPatientDoctors, patientUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Doctor
	for rows.Next() {
		var i Doctor
		if err := rows.Scan(
			&i.Uuid,
			&i.Name,
			&i.PhoneNumber,
			&i.Description,
			&i.Crp,
			&i.PixKey,
			&i.PaymentDetails,
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

const updateDoctor = `-- name: UpdateDoctor :one

UPDATE doctors
SET
    name = COALESCE(?1, name),
    description = COALESCE(
        ?2,
        description
    ),
    crp = COALESCE(?3, crp),
    pix_key = COALESCE(?4, pix_key),
    payment_details = COALESCE(
        ?5,
        payment_details
    ),
    updated_at = CURRENT_TIMESTAMP
WHERE
    uuid = ?6 RETURNING uuid, name, phone_number, description, crp, pix_key, payment_details, created_at, updated_at
`

type UpdateDoctorParams struct {
	Name           sql.NullString
	Description    sql.NullString
	Crp            sql.NullString
	PixKey         sql.NullString
	PaymentDetails sql.NullString
	Uuid           uuid.UUID
}

func (q *Queries) UpdateDoctor(ctx context.Context, arg UpdateDoctorParams) (Doctor, error) {
	row := q.db.QueryRowContext(ctx, updateDoctor,
		arg.Name,
		arg.Description,
		arg.Crp,
		arg.PixKey,
		arg.PaymentDetails,
		arg.Uuid,
	)
	var i Doctor
	err := row.Scan(
		&i.Uuid,
		&i.Name,
		&i.PhoneNumber,
		&i.Description,
		&i.Crp,
		&i.PixKey,
		&i.PaymentDetails,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
