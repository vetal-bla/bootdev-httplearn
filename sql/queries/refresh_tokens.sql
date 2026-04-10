-- name: CreateRefreshTokens :one
insert into refresh_tokens (
    token,
    created_at,
    updated_at,
    user_id,
    expires_at
)
values (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
returning *;

-- name: GetUserFromRefreshToken :one
select users.id as user_id, refresh_tokens.token, refresh_tokens.expires_at, refresh_tokens.revoked_at
from refresh_tokens
join users on refresh_tokens.user_id = users.id
where token = $1 limit 1;

-- name: RevokeRefreshToken :exec
update refresh_tokens
set updated_at = $2, revoked_at = $3
where token = $1;
