-- +goose Up
CREATE INDEX amount_idx ON transactions(id, amount);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX  amount_idx;