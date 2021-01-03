CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email varchar(20) NOT NULL,
    password varchar(256),
    UNIQUE (email)
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id SERIAL NOT NULL,
    text varchar(256),
    image_url text[],
    FOREIGN KEY (user_id) REFERENCES users(id)
);

