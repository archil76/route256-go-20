-- +goose Up
-- +goose StatementBegin
CREATE TABLE stocks (
                        id bigserial primary key,
                        total_count integer not null default 0,
                        reserved integer not null default 0
);

CREATE TABLE orders (
                        id bigserial primary key,
                        user_id bigint not null default 0
);

CREATE TABLE order_items (
                             id bigserial primary key,
                             order_id bigint not null,
                             user_id integer not null default 0,
                             status bpchar not null default ''
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stocks;
DROP TABLE orders;
DROP TABLE order_items;
-- +goose StatementEnd
