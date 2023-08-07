// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: doctor.sql

package db

import (
	"context"

	"github.com/gofrs/uuid"
)

const listPatientDoctors = `-- name: ListPatientDoctors :many

SELECT users.uuid, users.name, users.phone_number, users.description, users.crp, users.pix_key, users.payment_details, users.role, users.created_at, users.updated_at
FROM users
    JOIN patient_with_doctor ON users.uuid = patient_with_doctor.doctor_uuid
WHERE
    patient_with_doctor.patient_uuid = ?
`

func (q *Queries) ListPatientDoctors(ctx context.Context, patientUuid uuid.UUID) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listPatientDoctors, patientUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.Uuid,
			&i.Name,
			&i.PhoneNumber,
			&i.Description,
			&i.Crp,
			&i.PixKey,
			&i.PaymentDetails,
			&i.Role,
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
