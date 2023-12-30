-- +goose Up
-- +goose StatementBegin
CREATE table if not EXISTS gallerys(
    id serial primary key,
    title text,
    user_id int not null,
    desciption text 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table gallerys;
-- +goose StatementEnd
