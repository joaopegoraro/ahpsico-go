-- name: ListPatientDoctors :many

SELECT users.*
FROM users
    JOIN patient_with_doctor ON users.uuid = patient_with_doctor.doctor_uuid
WHERE
    patient_with_doctor.patient_uuid = ?;