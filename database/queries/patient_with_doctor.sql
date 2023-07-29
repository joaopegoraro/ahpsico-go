-- name: AddPatientDoctor :exec

INSERT INTO
    patient_with_doctor (doctor_uuid, patient_uuid)
VALUES (?, ?);

-- name: RemovePatientDoctor :exec

DELETE FROM
    patient_with_doctor
WHERE
    patient_uuid = ?
    AND doctor_uuid = ?;