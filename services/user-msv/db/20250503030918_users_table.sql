-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL CHECK (length(name) > 2), 
    registration_date TIMESTAMP NOT NULL DEFAULT NOW(),
    access_level INT NOT NULL CHECK (access_level IN (1, 2))
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS users;

-- +goose StatementEnd
