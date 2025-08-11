CREATE TABLE sessions(
    id SERIAL PRIMARY KEY,
    user_id IN UNIQUE,
    token_hash TEXT UNIQUE NOT NULL
);