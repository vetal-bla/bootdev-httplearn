-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllUsers :exec
delete from users;

-- name: GetUserByMail :one
select * from users
where email = $1;

-- name: UpdateEmailAndPawword :one
update users set email = $1, hashed_password = $2, updated_at = NOW()
where id = $3
returning *;
