-- +goose Up
-- +goose StatementBegin
CREATE TABLE passsword_resets (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    token_hash TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL
);
-- Alternative method to declare foreign key:
-- FOREIGN KEY (user_id) REFERENCES users (id)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE passsword_resets;
-- +goose StatementEnd
