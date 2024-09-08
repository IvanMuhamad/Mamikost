-- name: UploadImage :one
INSERT INTO rent_properties_images (frim_filename, frim_default, frim_repo_id) 
VALUES ($1, $2, $3) 
RETURNING frim_id;

-- name: GetAllImages :many
SELECT frim_id, frim_filename, frim_default, frim_repo_id 
FROM rent_properties_images;

-- name: GetImageByID :one
SELECT frim_id, frim_filename, frim_default, frim_repo_id 
FROM rent_properties_images 
WHERE frim_id = $1;

-- name: DeleteAllImagesForProperty :exec
DELETE FROM rent_properties_images 
WHERE frim_repo_id = $1;

-- name: DeleteImageByID :exec
DELETE FROM rent_properties_images 
WHERE frim_id = $1;