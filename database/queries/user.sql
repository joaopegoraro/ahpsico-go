-- name: GetUser :one

SELECT * FROM users WHERE uuid = ? LIMIT 1;

-- name: GetUserByRole :one

SELECT * FROM users WHERE uuid = ? AND role = ? LIMIT 1;

-- name: GetUserByPhoneNumber :one

SELECT * FROM users WHERE phone_number = ? LIMIT 1;

-- name: CreateUser :one

INSERT INTO
    users (
        uuid,
        name,
        phone_number,
        description,
        crp,
        pix_key,
        payment_details,
        role
    )
VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateUser :one

UPDATE users
SET
    name = COALESCE(sqlc.narg ('name'), name),
    description = COALESCE(
        sqlc.narg ('description'),
        description
    ),
    crp = COALESCE(sqlc.narg ('crp'), crp),
    pix_key = COALESCE(sqlc.narg ('pix_key'), pix_key),
    payment_details = COALESCE(
        sqlc.narg ('payment_details'),
        payment_details
    ),
    updated_at = CURRENT_TIMESTAMP
WHERE
    uuid = sqlc.arg('uuid') RETURNING *;