-- +goose Up
CREATE TABLE account (
    last_id varchar(20) not null,
    balance integer not null
);

-- +goose Down
DROP TABLE account;
