-- +goose Up
-- +goose StatementBegin
CREATE TABLE Comments (
    id UUID NOT NULL PRIMARY KEY,
    article_id UUID NOT NULL,
    owner_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    content TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Comments;
-- +goose StatementEnd
