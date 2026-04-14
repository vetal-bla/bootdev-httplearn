-- name: CreateChirps :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteChirp :exec
delete from chirps where id = $1;

-- name: GetAllChirps :many
select * from chirps
order by created_at asc;

-- name: GetChirp :one
select * from chirps
where id = $1;

-- name: DeleteUserChirp :exec
delete from chirps
where id=$1 and user_id=$2;

-- name: GetAllChirpsByUserID :many
select * from chirps
where user_id = $1
order by created_at asc;
