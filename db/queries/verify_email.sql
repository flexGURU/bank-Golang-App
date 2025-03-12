-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
    username,
    email,
    secret_code
) 
VALUES ($1, $2, $3) 
RETURNING *;


-- name: UpdateVerifyEmail :one
UPDATE verify_emails
SET is_used = TRUE
WHERE id = @id AND
secret_code = @secretcode AND
expires_at > NOW()
RETURNING *;

