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
        doctor_uuid TEXT NOT NULL,
        FOREIGN KEY (patient_uuid) REFERENCES patients (uuid) ON DELETE CASCADE,
        FOREIGN KEY (doctor_uuid) REFERENCES doctors (uuid) ON DELETE CASCADE
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE patients;

DROP TABLE patient_with_doctor;

-- +goose StatementEnd