-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    sessions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        patient_uuid TEXT NOT NULL,
        doctor_uuid TEXT NOT NULL,
        date DATETIME NOT NULL,
        group_index INTEGER NOT NULL DEFAULT 0,
        type INTEGER NOT NULL DEFAULT 0,
        status INTEGER NOT NULL DEFAULT 0,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME,
        FOREIGN KEY (patient_uuid) REFERENCES patients (uuid) ON DELETE CASCADE,
        FOREIGN KEY (doctor_uuid) REFERENCES doctors (uuid) ON DELETE CASCADE
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE sessions;

-- +goose StatementEnd