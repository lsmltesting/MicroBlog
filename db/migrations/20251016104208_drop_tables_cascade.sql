-- +goose Up
-- +goose StatementBegin


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists users cascade;
drop table if exists posts cascade;
drop table if exists likes cascade;
-- +goose StatementEnd
