-- +goose Up
-- +goose StatementBegin
CREATE TABLE Users (
    id UUID NOT NULL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    nick VARCHAR(50) NOT NULL,
    description TEXT,
    birthday TIMESTAMP WITHOUT TIME ZONE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Users;
-- +goose StatementEnd
