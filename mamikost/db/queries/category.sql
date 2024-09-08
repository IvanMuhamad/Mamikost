-- name: CreateCategory :one
INSERT INTO category (cate_name)
VALUES ($1)
RETURNING cate_id;

-- name: GetAllCategories :many
SELECT cate_id, cate_name
FROM category;

-- name: GetCategoryByID :one
SELECT cate_id, cate_name
FROM category
WHERE cate_id = $1;

-- name: UpdateCategory :exec
UPDATE category
SET cate_name = $1
WHERE cate_id = $2;

-- name: DeleteCategory :exec
DELETE FROM category
WHERE cate_id = $1;