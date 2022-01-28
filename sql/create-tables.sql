-- create table
CREATE TABLE slugs (
    id SERIAL PRIMARY KEY,
    extid INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    locale VARCHAR(255) NOT NULL,
    products TEXT NOT NULL,
    topics TEXT NOT NULL,
    summary TEXT NOT NULL
);
