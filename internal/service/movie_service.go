// Package service contains business logic and service layer for the application.
package service

import (
	"context"
	"fmt"

	"github.com/mexirica/chi-template/internal/db/repository"
	"github.com/mexirica/chi-template/internal/models"
	"github.com/mexirica/chi-template/internal/o11y"
)

type Service interface {
	Create(ctx context.Context, payload models.CreateMovieRequest) error
	GetById(ctx context.Context, id int) (*models.Movie, error)
	GetList(ctx context.Context, page, limit int) (*models.GetMovieList, error)
	Delete(ctx context.Context, id int) error
}

type MovieService struct {
	repo repository.MovieRepository
}

func NewMovieService(repo repository.MovieRepository) *MovieService {
	return &MovieService{
		repo: repo,
	}
}

func (s *MovieService) Create(ctx context.Context, payload models.CreateMovieRequest) error {
	ctx, span := o11y.Tracer().Start(ctx, "MovieService.Create")
	defer span.End()

	err := s.repo.Create(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to create movie: %w", err)
	}
	return nil
}

func (s *MovieService) GetById(ctx context.Context, id int) (*models.Movie, error) {
	ctx, span := o11y.Tracer().Start(ctx, "MovieService.GetById")
	defer span.End()

	movie, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie by id: %w", err)
	}
	return movie, nil
}

func (s *MovieService) GetList(ctx context.Context, page, limit int) (*models.GetMovieList, error) {
	ctx, span := o11y.Tracer().Start(ctx, "MovieService.GetList")
	defer span.End()

	movies, err := s.repo.GetList(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie list: %w", err)
	}
	return movies, nil
}

func (s *MovieService) Delete(ctx context.Context, id int) error {
	ctx, span := o11y.Tracer().Start(ctx, "MovieService.Delete")
	defer span.End()

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete movie: %w", err)
	}
	return nil
}
