-- +goose Up
-- +goose StatementBegin
CREATE TABLE notifications
(
    ID      SERIAL PRIMARY KEY,
    title   VARCHAR(255) NOT NULL,
    date    TIMESTAMP NOT NULL,
    user_id INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE notifications;
-- +goose StatementEnd
