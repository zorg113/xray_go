-- +goose Up
-- +goose StatementBegin
create table if not exists event (
    id serial PRIMARY KEY,
    title text not null,
    start_date timestamptz NOT NULL,
    end_date timestamptz NOT NULL,
    description text NOT NULL,
    owner_id integer NOT NULL,
    remind_in integer NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table event;
-- +goose StatementEnd
