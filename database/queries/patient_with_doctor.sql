-- name: AddPatientDoctor :exec

INSERT INTO
    patient_with_doctor (doctor_uuid, patient_uuid)
VALUES (?, ?);