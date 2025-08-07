-- name: CreateMovie :one
INSERT INTO movies (
    title, description, release_year, genre, director, rating
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetMovieByID :one
SELECT * FROM movies WHERE id = $1;

-- name: ListMovies :many
SELECT * FROM movies ORDER BY id DESC LIMIT $1 OFFSET $2;

-- name: UpdateMovie :one
UPDATE movies SET
    title = $2,
    description = $3,
    release_year = $4,
    genre = $5,
    director = $6,
    rating = $7
WHERE id = $1
RETURNING *;

-- name: DeleteMovie :exec
DELETE FROM movies WHERE id = $1;