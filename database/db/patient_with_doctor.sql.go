// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: patient_with_doctor.sql

package db

import (
	"context"

	"github.com/gofrs/uuid"
)

const addPatientDoctor = `-- name: AddPatientDoctor :exec

INSERT INTO
    patient_with_doctor (doctor_uuid, patient_uuid)
VALUES (?, ?)
`

type AddPatientDoctorParams struct {
	DoctorUuid  uuid.UUID
	PatientUuid uuid.UUID
}

func (q *Queries) AddPatientDoctor(ctx context.Context, arg AddPatientDoctorParams) error {
	_, err := q.db.ExecContext(ctx, addPatientDoctor, arg.DoctorUuid, arg.PatientUuid)
	return err
}

const removePatientDoctor = `-- name: RemovePatientDoctor :exec

DELETE FROM
    patient_with_doctor
WHERE
    patient_uuid = ?
    AND doctor_uuid = ?
`

type RemovePatientDoctorParams struct {
	PatientUuid uuid.UUID
	DoctorUuid  uuid.UUID
}

func (q *Queries) RemovePatientDoctor(ctx context.Context, arg RemovePatientDoctorParams) error {
	_, err := q.db.ExecContext(ctx, removePatientDoctor, arg.PatientUuid, arg.DoctorUuid)
	return err
}
