-- +goose Up
-- +goose StatementBegin
CREATE TABLE race_1234 (
    tag_id INT UNIQUE,
    start_time TIMESTAMP WITH TIME ZONE,
    end_time TIMESTAMP WITH TIME ZONE

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE race_1234;
-- +goose StatementEnd
