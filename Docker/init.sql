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

CREATE TABLE likes (
    user_id SERIAL NOT NULL,
    post_id SERIAL NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    CONSTRAINT user_like_post UNIQUE (user_id, post_id)
);

CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    user_id SERIAL NOT NULL,
    post_id SERIAL NOT NULL,
    text varchar(256),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
);

