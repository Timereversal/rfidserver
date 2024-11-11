-- +goose Up
-- +goose StatementBegin
CREATE TABLE runner(
    id SERIAL PRIMARY KEY ,
    name TEXT,
    lastName TEXT,
    birthDate date,
    distance integer,
    sex TEXT,
    category TEXT,
    tag_id integer,
    event_id integer,
    rank_category integer,
    rank_all integer,
    FOREIGN KEY (event_id) REFERENCES events (id),
    UNIQUE  (tag_id, event_id)

    
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
