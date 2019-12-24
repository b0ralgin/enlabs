-- +goose Up
CREATE TABLE transactions (
                      id varchar(50) NOT NULL,
                      amount integer not null,
                      state varchar(20) not null,
                      source varchar(20) not null,
                      PRIMARY KEY(id)
);

-- +goose Down
DROP table transactions;
