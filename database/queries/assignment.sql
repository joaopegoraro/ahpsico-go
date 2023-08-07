-- name: GetAssignment :one

SELECT assignments.* FROM assignments WHERE id = ? LIMIT 1;

-- name: ListPatientAssignments :many

SELECT
    assignments.id as assignment_id,
    assignments.title as assignment_title,
    assignments.description as assignment_description,
    assignments.status as assignment_status,
    assignments.patient_uuid as patient_uuid,
    sessions.id as session_id,
    sessions.date as session_date,
    doctors.uuid as doctor_uuid,
    doctors.name as doctor_name,
    doctors.description as doctor_description
FROM assignments
    JOIN users as doctors ON doctors.uuid = assignments.doctor_uuid
    JOIN sessions ON sessions.id = assignments.session_id
WHERE
    assignments.patient_uuid = ?;

-- name: ListDoctorPatientAssignments :many

SELECT
    assignments.id as assignment_id,
    assignments.title as assignment_title,
    assignments.description as assignment_description,
    assignments.status as assignment_status,
    assignments.patient_uuid as patient_uuid,
    sessions.id as session_id,
    sessions.date as session_date,
    doctors.uuid as doctor_uuid,
    doctors.name as doctor_name,
    doctors.description as doctor_description
FROM assignments
    JOIN users as doctors ON doctors.uuid = assignments.doctor_uuid
    JOIN sessions ON sessions.id = assignments.session_id
WHERE
    assignments.patient_uuid = ?
    AND assignments.doctor_uuid = ?;

-- name: ListPendingDoctorPatientAssignments :many

SELECT
    assignments.id as assignment_id,
    assignments.title as assignment_title,
    assignments.description as assignment_description,
    assignments.status as assignment_status,
    assignments.patient_uuid as patient_uuid,
    sessions.id as session_id,
    sessions.date as session_date,
    doctors.uuid as doctor_uuid,
    doctors.name as doctor_name,
    doctors.description as doctor_description
FROM assignments
    JOIN users as doctors ON doctors.uuid = assignments.doctor_uuid
    JOIN sessions ON sessions.id = assignments.session_id
WHERE
    assignments.patient_uuid = ?
    AND assignments.doctor_uuid = ?
    AND assignments.status = 0;

-- name: ListPendingPatientAssignments :many

SELECT
    assignments.id as assignment_id,
    assignments.title as assignment_title,
    assignments.description as assignment_description,
    assignments.status as assignment_status,
    assignments.patient_uuid as patient_uuid,
    sessions.id as session_id,
    sessions.date as session_date,
    doctors.uuid as doctor_uuid,
    doctors.name as doctor_name,
    doctors.description as doctor_description
FROM assignments
    JOIN users as doctors ON doctors.uuid = assignments.doctor_uuid
    JOIN sessions ON sessions.id = assignments.session_id
WHERE
    assignments.patient_uuid = ?
    AND assignments.status = 0;

-- name: CreateAssignment :one

INSERT INTO
    assignments (
        title,
        description,
        patient_uuid,
        doctor_uuid,
        session_id,
        status
    )
VALUES (?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateAssignment :one

UPDATE assignments
SET
    title = COALESCE(sqlc.narg ('title'), title),
    description = COALESCE(
        sqlc.narg ('description'),
        description
    ),
    status = COALESCE(sqlc.narg ('status'), status),
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = sqlc.arg('id') RETURNING *;

-- name: DeleteAssignment :exec

DELETE FROM assignments where id = ?;