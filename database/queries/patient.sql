-- name: GetPatient :one

SELECT * FROM patients WHERE uuid = ? LIMIT 1;

-- name: GetPatientByPhoneNumber :one

SELECT * FROM patients WHERE phone_number = ? LIMIT 1;

-- name: ListDoctorPatients :many

SELECT patients.*
FROM patients
    JOIN patient_with_doctor ON patients.uuid = patient_with_doctor.patient_uuid
WHERE
    patient_with_doctor.doctor_uuid = ?;

-- name: ListDoctorPatientsByPhoneNumber :many

SELECT patients.*
FROM patients
    JOIN patient_with_doctor ON patients.uuid = patient_with_doctor.patient_uuid
WHERE
    patient_with_doctor.doctor_uuid = ?
    AND patients.phone_number = ?;

-- name: CreatePatient :one

INSERT INTO
    patients (uuid, name, phone_number)
VALUES (?, ?, ?) RETURNING *;

-- name: UpdatePatient :one

UPDATE patients
SET
    name = COALESCE(sqlc.narg ('name'), name),
    updated_at = CURRENT_TIMESTAMP
WHERE
    uuid = sqlc.arg('uuid') RETURNING *;

-- name: DeletePatient :exec

DELETE FROM patients WHERE uuid = ?;