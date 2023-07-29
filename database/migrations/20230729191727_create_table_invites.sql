-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    invites (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        phone_number TEXT NOT NULL,
        patient_uuid TEXT NOT NULL,
        doctor_uuid TEXT NOT NULL,
        FOREIGN KEY (patient_uuid) REFERENCES patients (uuid) ON DELETE CASCADE,
        FOREIGN KEY (doctor_uuid) REFERENCES doctors (uuid) ON DELETE CASCADE
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE invites;

-- +goose StatementEnd