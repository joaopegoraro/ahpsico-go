-- name: GetSchedule :one

SELECT * FROM schedule WHERE id = ? LIMIT 1;

-- name: ListDoctorSchedule :many

SELECT * FROM schedule WHERE doctor_uuid = ?;

-- name: CreateSchedule :one

INSERT INTO schedule (doctor_uuid, date) VALUES (?, ?) RETURNING *;

-- name: DeleteSchedule :exec

DELETE FROM schedule WHERE id = ?;