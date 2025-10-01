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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stocks;
DROP TABLE orders;
DROP TABLE order_items;
-- +goose StatementEnd
