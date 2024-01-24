CREATE TABLE IF NOT EXISTS users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR NOT NULL,
                       surname VARCHAR NOT NULL,
                       patronymic VARCHAR,
                       age INT NOT NULL,
                       gender VARCHAR NOT NULL,
                       nationality VARCHAR NOT NULL
);
