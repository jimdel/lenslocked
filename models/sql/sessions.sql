CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    -- Inline method to declare foreign key 
    user_id INT UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    token_hash TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- Alternative method to declare foreign key:
    -- FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Add foreign key to existing table
-- ALTER TABLE sessions
-- ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id);

-- To delete all entries when deleting a parent entry:
-- ON DELETE CASCADE

-- Create an index
-- CREATE INDEX sessions_token_hash_idx ON sessions (token_hash, user_id, id);