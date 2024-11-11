-- +goose Up
-- +goose StatementBegin
CREATE TABLE categories(
    id SERIAL PRIMARY KEY ,
    name TEXT,
    age_low integer,
    age_high integer,
    participants integer,
    fk_event_id integer,
    FOREIGN KEY(fk_event_id) REFERENCES events (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE categories;
-- +goose StatementEnd
