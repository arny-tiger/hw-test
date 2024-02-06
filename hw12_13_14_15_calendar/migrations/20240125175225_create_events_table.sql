-- +goose Up
-- +goose StatementBegin
CREATE TABLE events
(
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    date        TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    duration INTERVAL NOT NULL,
    description TEXT,
    owner_id    INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
