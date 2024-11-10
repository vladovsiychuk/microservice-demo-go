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
    private_key text NOT NULL,
    public_key text NOT NULL,
    secondary_public_key text NOT NULL
);

CREATE TABLE session_tokens (
    id uuid PRIMARY KEY,
    email varchar(100),
    expires_at timestamptz NOT NULL
);