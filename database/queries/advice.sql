-- name: ListDoctorAdvices :many

SELECT * FROM advices WHERE doctor_uuid = ?;

-- name: ListPatientAdvicesFromDoctor :many

SELECT advices.*
FROM advices
    JOIN advice_with_patient ON advice.id = advice_with_patient.advice_id
WHERE
    advice_with_patient.patient_uuid = ?
    AND advices.doctor_uuid = ?;

-- name: CreateAdvice :exec

INSERT INTO advices (message, doctor_uuid) VALUES (?, ?);

-- name: DeleteAdvice :exec

DELETE FROM advices WHERE id = ?;