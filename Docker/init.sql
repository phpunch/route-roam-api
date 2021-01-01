CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email varchar(20),
    password varchar(256),
    UNIQUE (email)
);
