-- name: GetDoctor :one

SELECT * FROM doctors WHERE uuid = ? LIMIT 1;

-- name: ListPatientDoctors :many

SELECT doctors.*
FROM doctors
    JOIN patient_with_doctor ON doctors.uuid = patient_with_doctor.doctor_uuid
WHERE
    patient_with_doctor.patient_uuid = ?;

-- name: CreateDoctor :one

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
VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateDoctor :one

UPDATE doctors
SET
    name = COALESCE(sqlc.narg ('name'), title),
    description = COALESCE(
        sqlc.narg ('description'),
        title
    ),
    crp = COALESCE(sqlc.narg ('crp'), title),
    pix_key = COALESCE(sqlc.narg ('pix_key'), title),
    payment_details = COALESCE(
        sqlc.narg ('payment_details'),
        title
    ),
    updated_at = CURRENT_TIMESTAMP
WHERE
    uuid = sqlc.arg('uuid') RETURNING *;

-- name: DeleteDoctor :exec

DELETE FROM doctors WHERE uuid = ?;