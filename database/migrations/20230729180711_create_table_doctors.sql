-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    doctors (
        uuid TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        phone_number TEXT NOT NULL UNIQUE,
        description TEXT DEFAULT '' NOT NULL,
        crp TEXT DEFAULT '' NOT NULL,
        pix_key TEXT DEFAULT '' NOT NULL,
        payment_details TEXT DEFAULT '' NOT NULL,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE doctors;

-- +goose StatementEnd