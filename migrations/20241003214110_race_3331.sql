-- +goose Up
-- +goose StatementBegin
CREATE TABLE race_3331(
 tag_id INT UNIQUE,
 stage_0 TIMESTAMP WITH TIME ZONE,
 stage_1 TIMESTAMP WITH TIME ZONE,
 stage_2 TIMESTAMP WITH TIME ZONE,
 stage_3 TIMESTAMP WITH TIME ZONE,
 time_stage_1 INTERVAL GENERATED ALWAYS AS (AGE(stage_1,stage_0)) STORED
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE race_3331;
-- +goose StatementEnd
