-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS processed_events (
    event_id     text PRIMARY KEY,
    processed_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS post_like_counters (
    post_id    integer PRIMARY KEY,
    like_count bigint NOT NULL DEFAULT 0
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS post_like_counters;
DROP TABLE IF EXISTS processed_events;

-- +goose StatementEnd
