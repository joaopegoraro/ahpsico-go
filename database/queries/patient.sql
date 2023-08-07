-- name: GetDoctorPatientWithUuid :one

SELECT users.*
FROM users
    JOIN patient_with_doctor ON users.uuid = patient_with_doctor.patient_uuid
WHERE
    users.uuid = ?
    AND patient_with_doctor.doctor_uuid = ?
LIMIT 1;

-- name: ListDoctorPatients :many

SELECT users.*
FROM users
    JOIN patient_with_doctor ON users.uuid = patient_with_doctor.patient_uuid
WHERE
    patient_with_doctor.doctor_uuid = ?;

-- name: ListDoctorPatientsByPhoneNumber :many

SELECT users.*
FROM users
    JOIN patient_with_doctor ON users.uuid = patient_with_doctor.patient_uuid
WHERE
    patient_with_doctor.doctor_uuid = ?
    AND users.phone_number = ?;