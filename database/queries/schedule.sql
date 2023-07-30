-- name: ListDoctorSchedule :many

SELECT * FROM schedule WHERE doctor_uuid = ?;

-- name: CreateSchedule :exec

INSERT INTO schedule (doctor_uuid, date) VALUES (?, ?);

-- name: DeleteSchedule :exec

DELETE FROM schedule WHERE id = ?;