-- name: GetInvite :one

SELECT invites.* FROM invites WHERE invites.id = ? LIMIT 1;

-- name: GetDoctorInviteByPhoneNumber :one

SELECT invites.*
FROM invites
WHERE
    invites.doctor_uuid = ?
    AND invites.phone_number = ?
LIMIT 1;

-- name: ListDoctorInvites :many

SELECT
    invites.id as invite_id,
    invites.phone_number as invite_phone_number,
    invites.patient_uuid as invite_patient_uuid,
    invites.created_at as invite_created_at,
    doctors.uuid as doctor_uuid,
    doctors.name as doctor_name,
    doctors.description as doctor_description
FROM invites
    JOIN users as doctors on doctors.uuid = invites.doctor_uuid
WHERE invites.doctor_uuid = ?;

-- name: ListPatientInvites :many

SELECT
    invites.id as invite_id,
    invites.phone_number as invite_phone_number,
    invites.patient_uuid as invite_patient_uuid,
    invites.created_at as invite_created_at,
    doctors.uuid as doctor_uuid,
    doctors.name as doctor_name,
    doctors.description as doctor_description
FROM invites
    JOIN users as doctors on doctors.uuid = invites.doctor_uuid
WHERE invites.patient_uuid = ?;

-- name: CreateInvite :one

INSERT INTO
    invites (
        phone_number,
        doctor_uuid,
        patient_uuid
    )
VALUES (?, ?, ?) RETURNING *;

-- name: DeleteInvite :exec

DELETE FROM invites WHERE id = ?;