CREATE TABLE movies (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    release_year INT NOT NULL,
    genre VARCHAR(50)[], -- Tipo array de strings
    director VARCHAR(255),
    rating NUMERIC(3, 1)
);