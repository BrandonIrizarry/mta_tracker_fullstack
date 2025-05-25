-- +goose Up
CREATE TABLE routes(
       id TEXT PRIMARY KEY,
       long_name TEXT NOT NULL,
       short_name TEXT NOT NULL,
       description TEXT NOT NULL,
       select_bus_service BOOLEAN NOT NULL
);

-- +goose Down
DROP TABLE routes;
