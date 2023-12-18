Create TABLE sessions(
    id serial primary key ,
    token_hash text unique not null ,
    user_id int unique
);