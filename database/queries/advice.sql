-- name: GetAdvice :one

SELECT * FROM advices WHERE id = ? LIMIT 1;

-- name: ListDoctorAdvices :many

SELECT
    advices.id as advice_id,
    advices.message as advice_message,
    advices.created_at as advice_created_at,
    patients.uuid as patient_uuid,
    doctors.uuid as doctor_uuid,
    doctors.name as doctor_name
FROM advices
    JOIN advice_with_patient ON advices.id = advice_with_patient.advice_id
    JOIN users as doctors ON advices.doctor_uuid = doctors.uuid
    JOIN users as patients ON advice_with_patient.patient_uuid = patients.uuid
WHERE advices.doctor_uuid = ?;

-- name: ListPatientAdvices :many

SELECT
    advices.id as advice_id,
    advices.message as advice_message,
    advices.created_at as advice_created_at,
    patients.uuid as patient_uuid,
    doctors.uuid as doctor_uuid,
    doctors.name as doctor_name
FROM advices
    JOIN advice_with_patient ON advices.id = advice_with_patient.advice_id
    JOIN users as doctors ON advices.doctor_uuid = doctors.uuid
    JOIN users as patients ON advice_with_patient.patient_uuid = patients.uuid
WHERE
    advice_with_patient.patient_uuid = ?
GROUP BY advices.id;

-- name: ListDoctorPatientAdvices :many

SELECT
    advices.id as advice_id,
    advices.message as advice_message,
    advices.created_at as advice_created_at,
    patients.uuid as patient_uuid,
    doctors.uuid as doctor_uuid,
    doctors.name as doctor_name
FROM advices
    JOIN advice_with_patient ON advices.id = advice_with_patient.advice_id
    JOIN users as doctors ON advices.doctor_uuid = doctors.uuid
    JOIN users as patients ON advice_with_patient.patient_uuid = patients.uuid
WHERE
    advice_with_patient.patient_uuid = ?
    AND advices.doctor_uuid = ?;

-- name: CreateAdvice :one

INSERT INTO advices (message, doctor_uuid) VALUES (?, ?) RETURNING *;

-- name: DeleteAdvice :exec

DELETE FROM advices WHERE id = ?;