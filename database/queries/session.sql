-- name: GetSession :one

SELECT * FROM sessions WHERE id = ? LIMIT 1;

-- name: ListDoctorSessions :many

SELECT sessions.* FROM sessions WHERE doctor_uuid = ?;

-- name: ListDoctorSessionsWithinDate :many

SELECT sessions.*
FROM sessions
WHERE
    doctor_uuid = sqlc.arg ('doctor_uuid')
    AND date >= sqlc.arg ('start_of_date')
    AND date <= sqlc.arg ('end_of_date');

-- name: ListDoctorSessionsByExactDate :many

SELECT sessions.* FROM sessions WHERE doctor_uuid = ? AND date = ?;

-- name: ListPatientSessions :many

SELECT sessions.* FROM sessions WHERE patient_uuid = ?;

-- name: ListDoctorPatientSessions :many

SELECT sessions.*
FROM sessions
WHERE
    patient_uuid = ?
    AND doctor_uuid = ?;

-- name: ListUpcomingDoctorPatientSessions :many

SELECT sessions.*
FROM sessions
WHERE
    patient_uuid = ?
    AND doctor_uuid = ?
    AND date >= CURRENT_TIMESTAMP;

-- name: ListUpcomingPatientSessions :many

SELECT sessions.*
FROM sessions
WHERE
    patient_uuid = ?
    AND date >= CURRENT_TIMESTAMP;

-- name: CreateSession :one

INSERT INTO
    sessions (
        patient_uuid,
        doctor_uuid,
        date,
        group_index,
        type,
        status
    )
VALUES (?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateSession :one

UPDATE sessions
SET
    date = COALESCE(sqlc.narg ('date'), date),
    status = COALESCE(sqlc.narg ('status'), status),
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = sqlc.arg('id') RETURNING *;

-- name: DeleteSession :exec

DELETE FROM sessions where id = ?;