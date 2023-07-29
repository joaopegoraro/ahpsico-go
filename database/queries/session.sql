-- name: GetSession :one

SELECT * FROM sessions WHERE id = ? LIMIT 1;

-- name: ListDoctorSessions :many

SELECT sessions.* FROM sessions WHERE doctor_uuid = ?;

-- name: ListPatientSessions :many

SELECT sessions.* FROM sessions WHERE patient_uuid = ?;

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