-- +goose Up
-- +goose StatementBegin
create table if not exists password_reset(
    id serial primary key ,
    user_id int unique,
    token_hash text,
    expired_at timestamp not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table password_reset;
-- +goose StatementEnd
