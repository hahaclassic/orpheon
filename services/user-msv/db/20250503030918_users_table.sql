-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL CHECK (length(name) > 2), -- Имя должно быть хотя бы 3 символа
    registration_date TIMESTAMP NOT NULL DEFAULT NOW(), -- Дата регистрации по умолчанию
    birth_date DATE CHECK (birth_date < registration_date), -- Дата рождения не может быть в будущем
    access_level INT NOT NULL CHECK (access_level IN (1, 2)) -- Ограничиваем возможные роли
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS users;

-- +goose StatementEnd
