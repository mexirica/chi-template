// Package models defines the application's data structures and domain models.
package models

type Movie struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ReleaseYear int      `json:"release_year"`
	Genre       []string `json:"genre"`
	Director    string   `json:"director"`
	Rating      float64  `json:"rating"`
}

type GetMovieList struct {
	Movies []GetMovieResponse `json:"movies"`
}

type CreateMovieRequest struct {
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description" validate:"required"`
	ReleaseYear int      `json:"release_year" validate:"required"`
	Genre       []string `json:"genre" validate:"required"`
	Director    string   `json:"director" validate:"required"`
	Rating      float64  `json:"rating" validate:"required"`
}

type GetMovieResponse struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ReleaseYear int      `json:"release_year"`
	Genre       []string `json:"genre"`
	Director    string   `json:"director"`
	Rating      float64  `json:"rating"`
}
