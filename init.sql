CREATE DATABASE redhub;

\c redhub;

CREATE TABLE Users (
    id UUID NOT NULL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    nick VARCHAR(50) NOT NULL,
    description TEXT,
    birthday TIMESTAMP WITHOUT TIME ZONE
);

CREATE TABLE Articles (
    id UUID NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    owner_id UUID NOT NULL
);

CREATE TABLE Comments (
    id UUID NOT NULL PRIMARY KEY,
    article_id UUID NOT NULL,
    owner_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    content TEXT NOT NULL
);