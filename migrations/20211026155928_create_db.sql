-- +goose Up
CREATE TABLE IF NOT EXISTS metrics
(
    ID    TEXT NOT NULL,
    MType TEXT NOT NULL,
    Delta INTEGER,
    Value FLOAT,
    Hash  TEXT,
    PRIMARY KEY (ID)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
