CREATE TABLE contacts (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) NOT NULL,
    contact_id INT REFERENCES users(id) NOT NULL,
    UNIQUE(user_id, contact_id)
);