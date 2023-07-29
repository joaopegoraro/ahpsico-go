-- name: ListDoctorInvites :many

SELECT * FROM invites WHERE doctor_uuid = ?;

-- name: ListDoctorInvitesWithPhoneNumber :many

SELECT *
FROM invites
WHERE
    doctor_uuid = ?
    AND phone_number = ?;

-- name: ListPatientInvites :many

SELECT * FROM invites WHERE patient_uuid = ?;

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