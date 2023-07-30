-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    assignments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        description TEXT NOT NULL,
        patient_uuid TEXT NOT NULL,
        doctor_uuid TEXT NOT NULL,
        session_id INTEGER NOT NULL,
        status INTEGER NOT NULL DEFAULT 0,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME,
        FOREIGN KEY (patient_uuid) REFERENCES patients (uuid) ON DELETE CASCADE,
        FOREIGN KEY (doctor_uuid) REFERENCES doctors (uuid) ON DELETE CASCADE,
        FOREIGN KEY (session_id) REFERENCES sessions (id) ON DELETE CASCADE
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE assignments;

-- +goose StatementEnd