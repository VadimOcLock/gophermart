CREATE TABLE sessions (
                          id SERIAL PRIMARY KEY,
                          user_id INT REFERENCES users(id),
                          token TEXT NOT NULL,
                          expires_at TIMESTAMP NOT NULL
);