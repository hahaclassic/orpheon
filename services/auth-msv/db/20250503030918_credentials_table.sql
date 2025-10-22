-- +goose Up
-- +goose StatementBegin

CREATE TABLE credentials (
    user_id UUID PRIMARY KEY,
    login TEXT NOT NULL UNIQUE CHECK (length(login) > 2), -- Логин должен быть длиннее 3 символов
    password TEXT NOT NULL 
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS credentials;

-- +goose StatementEnd
