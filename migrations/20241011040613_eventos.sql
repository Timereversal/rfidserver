-- +goose Up
-- +goose StatementBegin
CREATE TABLE events(
    id SERIAL PRIMARY KEY ,
    name TEXT,
    location TEXT,
    time TIMESTAMP WITH TIME ZONE,
    categories JSONB,
--     subcategories text,
    distance integer[],
    status TEXT,
    event_picture TEXT, --link to file in filesystem
    event_video TEXT,
    event_sport TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
