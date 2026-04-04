-- +goose Up
-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE comments;
-- +goose StatementEnd
