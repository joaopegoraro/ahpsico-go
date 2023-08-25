-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    advices (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        message TEXT NOT NULL,
        doctor_uuid TEXT NOT NULL,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (doctor_uuid) REFERENCES users (uuid) ON DELETE CASCADE
    );

CREATE TABLE
    advice_with_patient (
        id INTEGER PRIMARY KEY,
        advice_id INTEGER NOT NULL,
        patient_uuid TEXT NOT NULL,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (advice_id) REFERENCES advices (id) ON DELETE CASCADE,
        FOREIGN KEY (patient_uuid) REFERENCES users (uuid) ON DELETE CASCADE
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE advices;

DROP TABLE advice_with_patient;

-- +goose StatementEnd