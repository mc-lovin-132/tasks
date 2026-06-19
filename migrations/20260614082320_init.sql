-- +goose Up
CREATE TABLE statuses(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE tasks(
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    deadline TIMESTAMP,
    status INTEGER NOT NULL REFERENCES statuses(id),
    creator_id INTEGER NOT NULL
);

-- +goose Down
DROP TABLE tasks;

DROP TABLE statuses;
