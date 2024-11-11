-- +goose Up
-- +goose StatementBegin
CREATE TABLE race_v1(
   tag_id INTEGER,
   event_id INTEGER,
   stage_0 TIMESTAMP WITH TIME ZONE,
   stage_1 TIMESTAMP WITH TIME ZONE,
   stage_2 TIMESTAMP WITH TIME ZONE,
   stage_3 TIMESTAMP WITH TIME ZONE,
   time_stage_1 INTERVAL GENERATED ALWAYS AS (AGE(stage_1,stage_0)) STORED,
   time_stage_2 INTERVAL GENERATED ALWAYS AS (AGE(stage_2,stage_0)) STORED,
   time_stage_3 INTERVAL GENERATED ALWAYS AS (AGE(stage_3,stage_0)) STORED
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE race_v1;
-- +goose StatementEnd
