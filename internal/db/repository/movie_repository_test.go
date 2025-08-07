// movie_repository_test.go
// Unit tests for the PsqlMovieRepository using GoMock.
package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	mock_repository "github.com/mexirica/chi-template/internal/db/repository"
	"github.com/mexirica/chi-template/internal/models"
)

func TestMovieRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockMovieRepository(ctrl)
	movie := models.CreateMovieRequest{
		Title:       "Test Movie",
		Description: "A test movie",
		ReleaseYear: 2024,
		Genre:       []string{"Action"},
		Director:    "John Doe",
		Rating:      8.5,
	}
	mockRepo.EXPECT().Create(gomock.Any(), movie).Return(nil)

	err := mockRepo.Create(context.Background(), movie)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestMovieRepository_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockMovieRepository(ctrl)
	movie := &models.Movie{ID: 1, Title: "Test Movie"}
	mockRepo.EXPECT().GetById(gomock.Any(), 1).Return(movie, nil)

	result, err := mockRepo.GetById(context.Background(), 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
}

func TestMovieRepository_GetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockMovieRepository(ctrl)
	movieList := &models.GetMovieList{Movies: []models.GetMovieResponse{{ID: 1, Title: "Test Movie"}}}
	mockRepo.EXPECT().GetList(gomock.Any(), 1, 10).Return(movieList, nil)

	result, err := mockRepo.GetList(context.Background(), 1, 10)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(result.Movies) != 1 {
		t.Errorf("expected 1 movie, got %d", len(result.Movies))
	}
}

func TestMovieRepository_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockMovieRepository(ctrl)
	mockRepo.EXPECT().Delete(gomock.Any(), 1).Return(nil)

	err := mockRepo.Delete(context.Background(), 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestMovieRepository_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockMovieRepository(ctrl)
	movie := models.CreateMovieRequest{Title: "Error Movie"}
	mockRepo.EXPECT().Create(gomock.Any(), movie).Return(errors.New("db error"))

	err := mockRepo.Create(context.Background(), movie)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
