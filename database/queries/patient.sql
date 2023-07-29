-- name: GetPatient :one

SELECT * FROM patients WHERE uuid = ? LIMIT 1;

-- name: ListDoctorPatients :many

SELECT patients.*
FROM patients
    JOIN patient_with_doctor ON patients.uuid = patient_with_doctor.patient_uuid
WHERE
    patient_with_doctor.doctor_uuid = ?;

-- name: CreatePatient :one

INSERT INTO
    patients (uuid, name, phone_number)
VALUES (?, ?, ?) RETURNING *;

-- name: UpdatePatient :one

UPDATE patients
SET
    name = COALESCE(sqlc.narg ('name'), title)
WHERE
    uuid = sqlc.arg('uuid') RETURNING *;

-- name: DeletePatient :exec

DELETE FROM patients WHERE uuid = ?;