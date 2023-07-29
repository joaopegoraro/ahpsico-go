-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    patients (
        uuid TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        phone_number TEXT NOT NULL UNIQUE,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME
    );

CREATE TABLE
    patient_with_doctor (
        patient_uuid TEXT NOT NULL,
        doctor_uuid TEXT NOT NULL,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (patient_uuid) REFERENCES patients (uuid) ON DELETE CASCADE,
        FOREIGN KEY (doctor_uuid) REFERENCES doctors (uuid) ON DELETE CASCADE
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE patients;

DROP TABLE patient_with_doctor;

-- +goose StatementEnd