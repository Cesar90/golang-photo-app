CREATE TABLE sessions(
    id SERIAL PRIMARY KEY,
    user_id IN UNIQUE REFERENCES users (id),
    token_hash TEXT UNIQUE NOT NULL
);