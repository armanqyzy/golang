CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       email TEXT NOT NULL UNIQUE,
                       name TEXT NOT NULL,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
