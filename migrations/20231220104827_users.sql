-- +goose Up
-- +goose StatementBegin
create table users(
    id serial primary key ,
    name text not null ,
    email text unique not null,
    password_hash text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table users;
-- +goose StatementEnd
