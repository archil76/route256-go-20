-- +goose Up
-- +goose StatementBegin
CREATE TABLE comments (
                        id bigserial not null primary key,
                        user_id bigint not null default 0,
                        sku integer not null,
                        comment text CHECK (length(comment) <= 250),
                        created_at timestamptz NOT NULL DEFAULT now()
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE comments;
-- +goose StatementEnd
