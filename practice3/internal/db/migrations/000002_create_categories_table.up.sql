CREATE TABLE categories (
                            id SERIAL PRIMARY KEY,
                            name TEXT NOT NULL,
                            user_id INTEGER,
                            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                            UNIQUE (user_id, name)
);

CREATE INDEX idx_categories_user_id ON categories(user_id);