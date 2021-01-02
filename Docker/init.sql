CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email varchar(20),
    password varchar(256),
    UNIQUE (email)
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id varchar(20),
    text varchar(256),
    image_url text[]
);
