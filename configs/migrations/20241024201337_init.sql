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

CREATE TABLE keys (
    private_key varchar(100) NOT NULL,
    public_key varchar(100) NOT NULL,
    secondary_public_key varchar(100) NOT NULL,
);