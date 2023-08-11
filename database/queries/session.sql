-- name: GetSession :one

SELECT * FROM sessions WHERE id = ? LIMIT 1;

-- name: GetSessionWithParticipants :one

SELECT
    sessions.id as session_id,
    sessions.date as session_date,
    sessions.group_index as session_group_index,
    sessions.type as session_type,
    sessions.status as session_status,
    sessions.payment_status as session_payment_status,
    sessions.created_at as session_created_at,
    doctors.uuid as doctor_uuid,
    doctors.name as doctor_name,
    doctors.description as doctor_description,
    patients.uuid as patient_uuid,
    patients.name as patient_name,
    patients.phone_number as patient_phone_number
FROM sessions
    JOIN users as doctors ON doctors.uuid = sessions.doctor_uuid
    JOIN users as patients ON patients.uuid = sessions.patient_uuid
WHERE sessions.id = ?
LIMIT 1;

-- name: GetDoctorSessionByExactDate :one

SELECT sessions.*
FROM sessions
WHERE doctor_uuid = ? AND date = ?
LIMIT 1;

-- name: ListDoctorSessions :many

SELECT
    sessions.id as session_id,
    sessions.date as session_date,
    sessions.group_index as session_group_index,
    sessions.type as session_type,
    sessions.status as session_status,
    sessions.payment_status as session_payment_status,
    sessions.created_at as session_created_at,
    doctors.uuid as doctor_uuid,
    doctors.name as doctor_name,
    doctors.description as doctor_description,
    patients.uuid as patient_uuid,
    patients.name as patient_name,
    patients.phone_number as patient_phone_number
FROM sessions
    JOIN users as doctors ON doctors.uuid = sessions.doctor_uuid
    JOIN users as patients ON patients.uuid = sessions.patient_uuid
WHERE sessions.doctor_uuid = ?;

-- name: ListDoctorActiveSessions :many

SELECT *
FROM sessions
WHERE
    doctor_uuid = ?
    AND status != 2
    AND status != 3;

-- name: ListDoctorSessionsWithinDate :many

SELECT
    sessions.id as session_id,
    sessions.date as session_date,
    sessions.group_index as session_group_index,
    sessions.type as session_type,
    sessions.status as session_status,
    sessions.payment_status as session_payment_status,
    sessions.created_at as session_created_at,
    doctors.uuid as doctor_uuid,
    doctors.name as doctor_name,
    doctors.description as doctor_description,
    patients.uuid as patient_uuid,
    patients.name as patient_name,
    patients.phone_number as patient_phone_number
FROM sessions
    JOIN users as doctors ON doctors.uuid = sessions.doctor_uuid
    JOIN users as patients ON patients.uuid = sessions.patient_uuid
WHERE
    sessions.doctor_uuid = sqlc.arg ('doctor_uuid')
    AND sessions.date >= sqlc.arg ('start_of_date')
    AND sessions.date <= sqlc.arg ('end_of_date');

-- name: ListPatientSessions :many

SELECT
    sessions.id as session_id,
    sessions.date as session_date,
    sessions.group_index as session_group_index,
    sessions.type as session_type,
    sessions.status as session_status,
    sessions.payment_status as session_payment_status,
    sessions.created_at as session_created_at,
    doctors.uuid as doctor_uuid,
    doctors.name as doctor_name,
    doctors.description as doctor_description,
    patients.uuid as patient_uuid,
    patients.name as patient_name,
    patients.phone_number as patient_phone_number
FROM sessions
    JOIN users as doctors ON doctors.uuid = sessions.doctor_uuid
    JOIN users as patients ON patients.uuid = sessions.patient_uuid
WHERE sessions.patient_uuid = ?;

-- name: ListDoctorPatientSessions :many

SELECT
    sessions.id as session_id,
    sessions.date as session_date,
    sessions.group_index as session_group_index,
    sessions.type as session_type,
    sessions.status as session_status,
    sessions.payment_status as session_payment_status,
    sessions.created_at as session_created_at,
    doctors.uuid as doctor_uuid,
    doctors.name as doctor_name,
    doctors.description as doctor_description,
    patients.uuid as patient_uuid,
    patients.name as patient_name,
    patients.phone_number as patient_phone_number
FROM sessions
    JOIN users as doctors ON doctors.uuid = sessions.doctor_uuid
    JOIN users as patients ON patients.uuid = sessions.patient_uuid
WHERE
    sessions.patient_uuid = ?
    AND sessions.doctor_uuid = ?;

-- name: ListUpcomingDoctorPatientSessions :many

SELECT
    sessions.id as session_id,
    sessions.date as session_date,
    sessions.group_index as session_group_index,
    sessions.type as session_type,
    sessions.status as session_status,
    sessions.payment_status as session_payment_status,
    sessions.created_at as session_created_at,
    doctors.uuid as doctor_uuid,
    doctors.name as doctor_name,
    doctors.description as doctor_description,
    patients.uuid as patient_uuid,
    patients.name as patient_name,
    patients.phone_number as patient_phone_number
FROM sessions
    JOIN users as doctors ON doctors.uuid = sessions.doctor_uuid
    JOIN users as patients ON patients.uuid = sessions.patient_uuid
WHERE
    sessions.patient_uuid = ?
    AND sessions.doctor_uuid = ?
    AND sessions.date >= CURRENT_TIMESTAMP;

-- name: ListUpcomingPatientSessions :many

SELECT
    sessions.id as session_id,
    sessions.date as session_date,
    sessions.group_index as session_group_index,
    sessions.type as session_type,
    sessions.status as session_status,
    sessions.payment_status as session_payment_status,
    sessions.created_at as session_created_at,
    doctors.uuid as doctor_uuid,
    doctors.name as doctor_name,
    doctors.description as doctor_description,
    patients.uuid as patient_uuid,
    patients.name as patient_name,
    patients.phone_number as patient_phone_number
FROM sessions
    JOIN users as doctors ON doctors.uuid = sessions.doctor_uuid
    JOIN users as patients ON patients.uuid = sessions.patient_uuid
WHERE
    sessions.patient_uuid = ?
    AND sessions.date >= CURRENT_TIMESTAMP;

-- name: CreateSession :one

INSERT INTO
    sessions (
        patient_uuid,
        doctor_uuid,
        date,
        group_index,
        type,
        status,
        payment_status
    )
VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id;

-- name: UpdateSession :one

UPDATE sessions
SET
    date = COALESCE(sqlc.narg ('date'), date),
    status = COALESCE(sqlc.narg ('status'), status),
    payment_status = COALESCE(
        sqlc.narg ('payment_status'),
        payment_status
    ),
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = sqlc.arg('id') RETURNING id;