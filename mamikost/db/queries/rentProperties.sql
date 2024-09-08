-- name: CreateRentProperty :one
INSERT INTO rent_properties (repo_name, repo_desc, repo_price, repo_cate_id)
VALUES ($1, $2, $3, $4)
RETURNING repo_id;

-- name: GetAllRentProperties :many
SELECT repo_id, repo_name, repo_desc, repo_price, repo_modified, repo_cate_id 
FROM rent_properties;

-- name: GetRentPropertyByID :one
SELECT repo_id, repo_name, repo_desc, repo_price, repo_modified, repo_cate_id 
FROM rent_properties 
WHERE repo_id = $1;

-- name: UpdateRentProperty :exec
UPDATE rent_properties
SET repo_name = $1, repo_desc = $2, repo_price = $3, repo_cate_id = $4, repo_modified = CURRENT_TIMESTAMP
WHERE repo_id = $5;

-- name: DeleteRentProperty :exec
DELETE FROM rent_properties
WHERE repo_id = $1;