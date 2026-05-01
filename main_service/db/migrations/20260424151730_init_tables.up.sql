CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name TEXT NOT NULL,
                       email TEXT UNIQUE NOT NULL,
                       password TEXT NOT NULL
);

CREATE TABLE boards (
                        id SERIAL PRIMARY KEY,
                        title TEXT NOT NULL,
                        user_id INT REFERENCES users(id)
);

CREATE TABLE tasks (
                       id SERIAL PRIMARY KEY,
                       title TEXT NOT NULL,
                       content TEXT,
                       status TEXT,
                       board_id INT REFERENCES boards(id)
);