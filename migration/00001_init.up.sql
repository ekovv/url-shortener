CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    Original text NOT NULL,
    Short text NOT NULL,
    cookie text NOT NULL
);
ALTER TABLE urls
    ADD CONSTRAINT unique_url
        UNIQUE (Original);