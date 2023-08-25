-- name: CreateAdviceWithPatient :exec

INSERT INTO
    advice_with_patient (advice_id, patient_uuid)
VALUES (?, ?);