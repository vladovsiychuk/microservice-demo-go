-- +goose up
CREATE TABLE posts (
    id uuid PRIMARY KEY,
    content text,
    is_private boolean
);

CREATE TABLE comments (
    id uuid PRIMARY KEY,
    post_id uuid NOT NULL,
    content text
);