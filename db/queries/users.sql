-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    full_name,
    email
) 
VALUES ($1, $2, $3, $4) 
RETURNING *;


-- name: GetUser :one
SELECT * FROM users 
WHERE username = $1 LIMIT 1;


-- name: GetUsers :many
SELECT username, full_name, email FROM users;

-- name: UpdateUserVerification :one
UPDATE users 
SET is_verified = TRUE
WHERE username = @username
RETURNING *;