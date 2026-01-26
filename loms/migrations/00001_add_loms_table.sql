-- +goose Up
-- +goose StatementBegin
CREATE TABLE stocks (
                        id bigserial primary key,
                        total_count integer not null default 0,
                        reserved integer not null default 0
);

CREATE TABLE orders (
                        id bigserial primary key,
                        user_id bigint not null default 0,
                        status bpchar not null default ''
);

CREATE TABLE order_items (
                        order_id bigint REFERENCES orders(id) ON DELETE CASCADE,
                        sku integer not null,
                        count integer not null default 0
);

CREATE TABLE if not exists outbox
(
    id         bigserial PRIMARY KEY,
    key        text,
    payload    jsonb       NOT NULL,
    status     text        NOT NULL DEFAULT 'new', -- new | sent | error
    created_at timestamptz NOT NULL DEFAULT now(),
    sent_at    timestamptz
);
CREATE INDEX ON outbox (status, created_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stocks;
DROP TABLE order_items;
DROP TABLE orders;
DROP TABLE outbox;
-- +goose StatementEnd
