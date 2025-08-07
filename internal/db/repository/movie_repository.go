// Package repository provides database repository implementations for the application.
package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mexirica/chi-template/internal/db/sqlc"
	"github.com/mexirica/chi-template/internal/models"
	"github.com/mexirica/chi-template/internal/o11y"
)

type MovieRepository interface {
	Create(ctx context.Context, movie models.CreateMovieRequest) error
	GetById(ctx context.Context, id int) (*models.Movie, error)
	GetList(ctx context.Context, page, limit int) (*models.GetMovieList, error)
	Delete(ctx context.Context, id int) error
}

type PsqlMovieRepository struct {
	sqlc.Querier
}

func NewMovieRepository(conn *pgxpool.Pool) *PsqlMovieRepository {
	return &PsqlMovieRepository{
		Querier: sqlc.New(conn),
	}
}

func (r *PsqlMovieRepository) Create(ctx context.Context, movie models.CreateMovieRequest) error {
	ctx, span := o11y.Tracer().Start(ctx, "PsqlMovieRepository.Create")
	defer span.End()
	_, err := r.CreateMovie(ctx, sqlc.CreateMovieParams{
		Title:       movie.Title,
		Description: pgtype.Text{String: movie.Description, Valid: true},
		ReleaseYear: int32(movie.ReleaseYear),
		Genre:       movie.Genre,
		Director:    pgtype.Text{String: movie.Director, Valid: true},
		Rating:      float64ToPgNumeric(movie.Rating),
	})
	return err
}

func (r *PsqlMovieRepository) GetById(ctx context.Context, id int) (*models.Movie, error) {
	ctx, span := o11y.Tracer().Start(ctx, "PsqlMovieRepository.GetById")
	defer span.End()
	movie, err := r.GetMovieByID(ctx, int64(id))
	if err != nil {
		return nil, err
	}

	rating, err := movie.Rating.Float64Value()
	if err != nil {
		return nil, fmt.Errorf("failed to get movie rating: %w", err)
	}
	return &models.Movie{
		ID:          movie.ID,
		Title:       movie.Title,
		Description: movie.Description.String,
		ReleaseYear: int(movie.ReleaseYear),
		Genre:       movie.Genre,
		Director:    movie.Director.String,
		Rating:      rating.Float64,
	}, nil
}

func (r *PsqlMovieRepository) GetList(ctx context.Context, page, limit int) (*models.GetMovieList, error) {
	ctx, span := o11y.Tracer().Start(ctx, "PsqlMovieRepository.GetList")
	defer span.End()

	offset := (page - 1) * limit
	movies, err := r.ListMovies(ctx, sqlc.ListMoviesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]models.GetMovieResponse, 0, len(movies))
	for _, m := range movies {
		rating, _ := m.Rating.Float64Value()
		result = append(result, models.GetMovieResponse{
			ID:          m.ID,
			Title:       m.Title,
			Description: m.Description.String,
			ReleaseYear: int(m.ReleaseYear),
			Genre:       m.Genre,
			Director:    m.Director.String,
			Rating:      rating.Float64,
		})
	}

	return &models.GetMovieList{
		Movies: result,
	}, nil
}

func (r *PsqlMovieRepository) Delete(ctx context.Context, id int) error {
	ctx, span := o11y.Tracer().Start(ctx, "PsqlMovieRepository.Delete")
	defer span.End()
	err := r.DeleteMovie(ctx, int64(id))
	if err != nil {
		return fmt.Errorf("failed to delete movie with id %d: %w", id, err)
	}
	return nil
}

// float64ToPgNumeric converts a float64 to pgtype.Numeric.
func float64ToPgNumeric(val float64) pgtype.Numeric {
	var num pgtype.Numeric
	_ = num.Scan(val)
	return num
}
