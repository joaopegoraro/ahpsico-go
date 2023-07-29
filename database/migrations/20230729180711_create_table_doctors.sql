-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    doctors (
        uuid TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        phone_number TEXT NOT NULL UNIQUE,
        description TEXT DEFAULT '' NOT NULL UNIQUE,
        crp TEXT DEFAULT '' NOT NULL UNIQUE,
        pix_key TEXT DEFAULT '' NOT NULL UNIQUE,
        payment_details TEXT DEFAULT '' NOT NULL UNIQUE
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE doctors;

-- +goose StatementEnd