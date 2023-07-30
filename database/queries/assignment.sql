-- name: ListPatientAssignments :many

SELECT assignments.* FROM assignments WHERE patient_uuid = ?;

-- name: ListDoctorPatientAssignments :many

SELECT assignments.*
FROM assignments
WHERE
    patient_uuid = ?
    AND doctor_uuid = ?;

-- name: ListPendingDoctorPatientAssignments :many

SELECT assignments.*
FROM assignments
WHERE
    patient_uuid = ?
    AND doctor_uuid = ?
    AND status = 0;

-- name: ListPendingPatientAssignments :many

SELECT assignments.*
FROM assignments
WHERE
    patient_uuid = ?
    AND status = 0;

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