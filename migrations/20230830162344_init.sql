-- +goose Up
-- +goose StatementBegin
CREATE TABLE apods
(
    title       TEXT NOT NULL,
    explanation TEXT NOT NULL,
    date        DATE PRIMARY KEY,
    url         TEXT NOT NULL,
    image_b64   TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE apods;
-- +goose StatementEnd
