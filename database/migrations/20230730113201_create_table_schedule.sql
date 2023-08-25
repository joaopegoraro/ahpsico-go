-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    schedule (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        doctor_uuid TEXT NOT NULL,
        date DATETIME NOT NULL,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME,
        FOREIGN KEY (doctor_uuid) REFERENCES users (uuid) ON DELETE CASCADE
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE schedule;

-- +goose StatementEnd