-- +goose Up
-- +goose StatementBegin
CREATE TABLE Articles (
    id UUID NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    tag VARCHAR(50) NOT NULL,
    owner_id UUID NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Articles;
-- +goose StatementEnd
