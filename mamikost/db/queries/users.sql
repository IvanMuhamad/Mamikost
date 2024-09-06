-- name: FindUserByUserPassword :one
select user_id,user_name,user_password,user_phone,user_token 
from users  
WHERE user_name = $1 and user_password = crypt($2, user_password);

-- name: FindUserByUsername :one
select user_id,user_name,user_password,user_phone,user_token 
from users  
WHERE user_name = $1;

-- name: FindUserByPhone :one
select user_id,user_name,user_password,user_phone,user_token 
from users  
WHERE user_phone  = $1;

-- name: CreateUser :one
INSERT INTO users(user_name,user_password,user_email, user_phone)
	VALUES
	($1, crypt($2,gen_salt('bf')),$3, $4)
    RETURNING *;

-- name: UpdateUserName :one
UPDATE users SET user_name = $1 WHERE user_id = $2
RETURNING *;

-- name: UpdateUserPhone :one
UPDATE users SET user_phone = $1 WHERE user_id = $2
RETURNING *;

-- name: UpdateToken :one
UPDATE users SET user_token = $1 WHERE user_id = $2
RETURNING *;

-- name: DeleteToken :exec
UPDATE users SET user_token = '' WHERE user_token = $1
RETURNING *;