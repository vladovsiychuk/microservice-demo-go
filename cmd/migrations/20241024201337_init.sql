-- +goose up
CREATE TABLE posts (
    id uuid PRIMARY KEY,
    content text,
    is_private boolean
);