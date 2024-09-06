CREATE TABLE category (
    cate_id SERIAL PRIMARY KEY,
    cate_name VARCHAR(25) UNIQUE NOT NULL
);

CREATE TABLE rent_properties (
    repo_id SERIAL PRIMARY KEY,
    repo_name VARCHAR(55) NOT NULL,
    repo_desc VARCHAR(125),
    repo_price DOUBLE PRECISION NOT NULL,
    repo_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    repo_cate_id INT REFERENCES category(cate_id)
);

CREATE TABLE rent_properties_images (
    frim_id SERIAL PRIMARY KEY,
    frim_filename VARCHAR(125) NOT NULL,
    frim_default CHAR(1),
    frim_repo_id INT REFERENCES rent_properties(repo_id)
);

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    user_name VARCHAR(25) UNIQUE NOT NULL,
    user_password VARCHAR(85) NOT NULL,
    user_email VARCHAR(25) UNIQUE NOT NULL,
    user_phone VARCHAR(15) UNIQUE NOT NULL,
    user_token VARCHAR(255)
);


CREATE TABLE order_rent_properties (
    orpo_id SERIAL PRIMARY KEY,
    orpo_purchase_no VARCHAR(25) UNIQUE NOT NULL,
    orpo_tax DOUBLE PRECISION,
    orpo_subtotal DOUBLE PRECISION,
    orpo_patrx_no VARCHAR(55),
    orpo_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    orpo_user_id INT REFERENCES users(user_id)
);

CREATE TABLE order_rent_properties_detail (
    orpd_id SERIAL PRIMARY KEY,
    orpd_qty_unit INT NOT NULL,
    orpd_price DOUBLE PRECISION NOT NULL,
    orpd_total_price DOUBLE PRECISION,
    orpd_orpo_id INT REFERENCES order_rent_properties(orpo_id),
    orpd_repo_id INT REFERENCES rent_properties(repo_id)
);

CREATE TABLE carts (
    cart_id SERIAL PRIMARY KEY,
    cart_user_id INT UNIQUE REFERENCES users(user_id),
    cart_fr_id INT UNIQUE REFERENCES rent_properties_images(frim_id),
    cart_start_date TIMESTAMP NOT NULL,
    cart_end_date DATE NOT NULL,
    cart_qty INT NOT NULL,
    cart_price DOUBLE PRECISION NOT NULL,
    cart_total_price DOUBLE PRECISION,
    cart_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    cart_status VARCHAR(15),
    cart_cart_id INT REFERENCES carts(cart_id)
);