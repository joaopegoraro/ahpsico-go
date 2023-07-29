-- name: ListDoctorAdvices :many

SELECT * FROM advices WHERE doctor_uuid = ?;

-- name: ListPatientAdvices :many

SELECT * FROM advices WHERE doctor_uuid = ?;

-- name: CreateAdvice :exec

INSERT INTO advices (message, doctor_uuid) VALUES (?, ?);

-- name: DeleteAdvice :exec

DELETE FROM advices WHERE id = ?;