
-- name: AddRoute :exec
INSERT INTO routes (
       id, long_name, short_name, description, select_bus_service
) VALUES (
  ?, ?, ?, ?, ?
);

-- name: GetAllRoutes :many
SELECT * FROM routes;

-- name: TestRouteTablePopulated :one
SELECT EXISTS (
       SELECT 1 FROM routes
);
