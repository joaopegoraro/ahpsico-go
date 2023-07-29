-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    patients (
        uuid TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        phone_number TEXT NOT NULL UNIQUE
    );

CREATE TABLE
    patient_with_doctor (
        id INTEGER PRIMARY KEY,
        patient_uuid TEXT NOT NULL,
        doctor_uuid TEXT NOT NULL
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE patients;

DROP TABLE patient_with_doctor;

-- +goose StatementEnd