CREATE TABLE expenses (
                          id SERIAL PRIMARY KEY,
                          user_id INTEGER NOT NULL,
                          category_id INTEGER NOT NULL,
                          amount NUMERIC(10, 2) NOT NULL CHECK (amount > 0),
                          currency CHAR(3) NOT NULL,
                          spent_at TIMESTAMP NOT NULL,
                          created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                          note TEXT,
                          FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                          FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
);

CREATE INDEX idx_expenses_user_id ON expenses(user_id);
CREATE INDEX idx_expenses_user_spent ON expenses(user_id, spent_at);