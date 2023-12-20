-- +goose Up
-- +goose StatementBegin
Create TABLE sessions(
    id serial primary key ,
    token_hash text unique not null ,
    user_id int unique
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Drop table sessions;
-- +goose StatementEnd
