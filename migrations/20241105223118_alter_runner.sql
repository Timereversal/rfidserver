-- +goose Up
-- +goose StatementBegin
ALTER TABLE runner ALTER column event_id SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE runner ALTER COLUMN event_id DROP NOT NULL ;
-- +goose StatementEnd
