CREATE TABLE users (
    Id SERIAL PRIMARY KEY,
    Firstname varchar(255),
    Lastname varchar(255) NOT NULL,
    Email varchar(255) NOT NULL
);