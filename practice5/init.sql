CREATE TABLE IF NOT EXISTS categories (
                                          id SERIAL PRIMARY KEY,
                                          name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS products (
                                        id SERIAL PRIMARY KEY,
                                        name TEXT NOT NULL,
                                        category_id INTEGER NOT NULL REFERENCES categories(id),
    price INTEGER NOT NULL
    );

INSERT INTO categories (name) VALUES
                                  ('phones'),
                                  ('laptops'),
                                  ('accessories'),
                                  ('tablets')
    ON CONFLICT (name) DO NOTHING;

INSERT INTO products (name, category_id, price) VALUES
                                                    ('iPhone 15', (SELECT id FROM categories WHERE name = 'phones'), 400000),
                                                    ('Samsung Galaxy S24', (SELECT id FROM categories WHERE name = 'phones'), 350000),
                                                    ('MacBook Pro', (SELECT id FROM categories WHERE name = 'laptops'), 800000),
                                                    ('Dell XPS 15', (SELECT id FROM categories WHERE name = 'laptops'), 600000),
                                                    ('iPad Pro', (SELECT id FROM categories WHERE name = 'tablets'), 300000),
                                                    ('AirPods Pro', (SELECT id FROM categories WHERE name = 'accessories'), 100000),
                                                    ('Magic Mouse', (SELECT id FROM categories WHERE name = 'accessories'), 35000);
