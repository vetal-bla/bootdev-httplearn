-- +goose Up
create table refresh_tokens (
    token text not null primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    user_id uuid not null,
    expires_at timestamp not null,
    revoked_at timestamp default null,
    constraint fk_user_id foreign key (user_id) references users(id) on delete cascade
);

-- +goose Down
drop table refresh_tokens;
