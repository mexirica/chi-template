// movie_service_test.go
// Unit tests for the MovieService using GoMock.
package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mexirica/chi-template/internal/models"
	"github.com/mexirica/chi-template/internal/service"
	mock_service "github.com/mexirica/chi-template/internal/service"
)

func TestMovieService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_service.NewMockService(ctrl)
	svc := service.NewMovieService(mockRepo)

	payload := models.CreateMovieRequest{
		Title:       "Test Movie",
		Description: "A test movie",
		ReleaseYear: 2024,
		Genre:       []string{"Action"},
		Director:    "John Doe",
		Rating:      8.5,
	}

	mockRepo.EXPECT().Create(gomock.Any(), payload).Return(nil)

	err := svc.Create(context.Background(), payload)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestMovieService_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_service.NewMockService(ctrl)
	svc := service.NewMovieService(mockRepo)

	movie := &models.Movie{ID: 1, Title: "Test Movie"}
	mockRepo.EXPECT().GetById(gomock.Any(), 1).Return(movie, nil)

	result, err := svc.GetById(context.Background(), 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
}

func TestMovieService_GetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_service.NewMockService(ctrl)
	svc := service.NewMovieService(mockRepo)

	movieList := &models.GetMovieList{Movies: []models.GetMovieResponse{{ID: 1, Title: "Test Movie"}}}
	mockRepo.EXPECT().GetList(gomock.Any(), 1, 10).Return(movieList, nil)

	result, err := svc.GetList(context.Background(), 1, 10)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(result.Movies) != 1 {
		t.Errorf("expected 1 movie, got %d", len(result.Movies))
	}
}

func TestMovieService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_service.NewMockService(ctrl)
	svc := service.NewMovieService(mockRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), 1).Return(nil)

	err := svc.Delete(context.Background(), 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestMovieService_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_service.NewMockService(ctrl)
	svc := service.NewMovieService(mockRepo)

	payload := models.CreateMovieRequest{Title: "Error Movie"}
	mockRepo.EXPECT().Create(gomock.Any(), payload).Return(errors.New("db error"))

	err := svc.Create(context.Background(), payload)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
