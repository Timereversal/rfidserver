-- +goose Up
-- +goose StatementBegin
ALTER TABLE events ALTER COLUMN time SET DATA TYPE DATE  ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events ALTER COLUMN time SET DATA TYPE  TIMESTAMP WITH TIME ZONE;
-- +goose StatementEnd






