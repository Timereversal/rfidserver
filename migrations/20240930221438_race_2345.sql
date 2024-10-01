-- +goose Up
-- +goose StatementBegin
CREATE TABLE race_2345(
   tag_id INT UNIQUE,
   stage_1 TIMESTAMP WITH TIME ZONE,
   stage_2 TIMESTAMP WITH TIME ZONE,
   stage_3 TIMESTAMP WITH TIME ZONE,
   stage_4 TIMESTAMP WITH TIME ZONE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE race_2345;
-- +goose StatementEnd
