-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    invites (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        phone_number TEXT NOT NULL,
        patient_uuid TEXT NOT NULL,
        doctor_uuid TEXT NOT NULL,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (patient_uuid) REFERENCES users (uuid) ON DELETE CASCADE,
        FOREIGN KEY (doctor_uuid) REFERENCES users (uuid) ON DELETE CASCADE
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE invites;

-- +goose StatementEnd