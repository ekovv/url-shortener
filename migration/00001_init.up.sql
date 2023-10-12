CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    original text NOT NULL,
    short text NOT NULL,
    cookie text NOT NULL
);
ALTER TABLE urls
    ADD CONSTRAINT unique_url
        UNIQUE (Original);