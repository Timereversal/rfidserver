-- +goose Up
-- +goose StatementBegin
ALTER TABLE events ALTER column name set not null;
ALTER TABLE events ADD CONSTRAINT name_unique UNIQUE (name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events DROP CONSTRAINT IF EXISTS name_unique;
-- +goose StatementEnd
