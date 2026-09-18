-- +goose up
CREATE TABLE tags (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL
);

-- +goose down
DROP TABLE tags;

