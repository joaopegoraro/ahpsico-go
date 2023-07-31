-- name: GetSession :one

SELECT * FROM sessions WHERE id = ? LIMIT 1;

-- name: GetSessionWithParticipants :one

SELECT
    s.id as session_id,
    s.date as session_date,
    s.group_index as session_group_index,
    s.type as session_type,
    s.status as session_status,
    s.created_at as session_created_at,
    d.uuid as doctor_uuid,
    d.name as doctor_name,
    d.description as doctor_description,
    p.uuid as patient_uuid,
    p.name as patient_name,
    p.phone_number as patient_phone_number
FROM sessions s
    JOIN doctors d ON doctors.uuid = sessions.doctor_uuid
    JOIN patients p ON patients.uuid = sessions.patient_uuid
WHERE id = ?
LIMIT 1;

-- name: GetDoctorSessionByExactDate :one

SELECT sessions.*
FROM sessions
WHERE doctor_uuid = ? AND date = ?
LIMIT 1;

-- name: ListDoctorSessions :many

SELECT
    s.id as session_id,
    s.date as session_date,
    s.group_index as session_group_index,
    s.type as session_type,
    s.status as session_status,
    s.created_at as session_created_at,
    p.uuid as patient_uuid,
    p.name as patient_name,
    p.phone_number as patient_phone_number
FROM sessions s
    JOIN patients p ON patients.uuid = sessions.patient_uuid
WHERE s.doctor_uuid = ?;

-- name: ListDoctorActiveSessions :many

SELECT *
FROM sessions
WHERE
    doctor_uuid = ?
    AND status != 2
    AND status != 3;

-- name: ListDoctorSessionsWithinDate :many

SELECT
    s.id as session_id,
    s.date as session_date,
    s.group_index as session_group_index,
    s.type as session_type,
    s.status as session_status,
    s.created_at as session_created_at,
    p.uuid as patient_uuid,
    p.name as patient_name,
    p.phone_number as patient_phone_number
FROM sessions s
    JOIN patients p ON p.uuid = s.patient_uuid
WHERE
    s.doctor_uuid = sqlc.arg ('doctor_uuid')
    AND s.date >= sqlc.arg ('start_of_date')
    AND s.date <= sqlc.arg ('end_of_date');

-- name: ListPatientSessions :many

SELECT
    s.id as session_id,
    s.date as session_date,
    s.group_index as session_group_index,
    s.type as session_type,
    s.status as session_status,
    s.created_at as session_created_at,
    d.uuid as doctor_uuid,
    d.name as doctor_name,
    d.description as doctor_description,
    p.uuid as patient_uuid,
    p.name as patient_name,
    p.phone_number as patient_phone_number
FROM sessions s
    JOIN doctors d ON d.uuid = s.doctor_uuid
    JOIN patients p ON p.uuid = s.patient_uuid
WHERE s.patient_uuid = ?;

-- name: ListDoctorPatientSessions :many

SELECT
    s.id as session_id,
    s.date as session_date,
    s.group_index as session_group_index,
    s.type as session_type,
    s.status as session_status,
    s.created_at as session_created_at,
    d.uuid as doctor_uuid,
    d.name as doctor_name,
    d.description as doctor_description,
    p.uuid as patient_uuid,
    p.name as patient_name,
    p.phone_number as patient_phone_number
FROM sessions s
    JOIN doctors d ON d.uuid = s.doctor_uuid
    JOIN patients p ON p.uuid = s.patient_uuid
WHERE
    s.patient_uuid = ?
    AND s.doctor_uuid = ?;

-- name: ListUpcomingDoctorPatientSessions :many

SELECT
    s.id as session_id,
    s.date as session_date,
    s.group_index as session_group_index,
    s.type as session_type,
    s.status as session_status,
    s.created_at as session_created_at,
    d.uuid as doctor_uuid,
    d.name as doctor_name,
    d.description as doctor_description,
    p.uuid as patient_uuid,
    p.name as patient_name,
    p.phone_number as patient_phone_number
FROM sessions s
    JOIN doctors d ON d.uuid = s.doctor_uuid
    JOIN patients p ON p.uuid = s.patient_uuid
WHERE
    s.patient_uuid = ?
    AND s.doctor_uuid = ?
    AND s.date >= CURRENT_TIMESTAMP;

-- name: ListUpcomingPatientSessions :many

SELECT
    s.id as session_id,
    s.date as session_date,
    s.group_index as session_group_index,
    s.type as session_type,
    s.status as session_status,
    s.created_at as session_created_at,
    d.uuid as doctor_uuid,
    d.name as doctor_name,
    d.description as doctor_description,
    p.uuid as patient_uuid,
    p.name as patient_name,
    p.phone_number as patient_phone_number
FROM sessions s
    JOIN doctors d ON d.uuid = s.doctor_uuid
    JOIN patients p ON p.uuid = s.patient_uuid
WHERE
    s.patient_uuid = ?
    AND s.date >= CURRENT_TIMESTAMP;

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