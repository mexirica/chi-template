// Package handler provides HTTP handlers for movie-related endpoints.
package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mexirica/chi-template/internal/helpers"
	"github.com/mexirica/chi-template/internal/models"
	"github.com/mexirica/chi-template/internal/o11y"
	"github.com/mexirica/chi-template/internal/service"
	"github.com/mexirica/chi-template/internal/validation"
)

type MovieHandler struct {
	s service.Service
}

func NewMovieHandler(service service.Service) *MovieHandler {
	return &MovieHandler{
		s: service,
	}
}

// CreateMovie godoc
// @Summary Create a new movie
// @Description Create a new movie with the provided details
// @Tags movies
// @Accept json
// @Produce json
// @Param movie body models.CreateMovieRequest true "Movie to create"
// @Success 201 {object} models.GetMovieResponse
// @Failure 400 {object} types.JsonResponse
// @Router /movies [post]
func (h *MovieHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, span := o11y.Tracer().Start(r.Context(), "MovieHandler.Create")
	defer span.End()
	var payload models.CreateMovieRequest

	errList, err := validation.BindAndValidate(r, payload)
	if err != nil {
		if errList != nil {
			helpers.WriteJSON(w, http.StatusBadRequest, errList)
			return
		}
	}

	err = h.s.Create(ctx, payload)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusCreated, nil)
}

// GetMovie godoc
// @Summary Get a movie by ID
// @Description Get details of a movie by its ID
// @Tags movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} models.GetMovieResponse
// @Failure 404 {object} types.JsonResponse
// @Router /movies/{id} [get]
func (h *MovieHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx, span := o11y.Tracer().Start(r.Context(), "MovieHandler.GetById")
	defer span.End()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	movie, err := h.s.GetById(ctx, id)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if movie == nil {
		helpers.WriteJSON(w, http.StatusNotFound, nil)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, movie)
}

// ListMovies godoc
// @Summary List movies
// @Description Get a paginated list of movies
// @Tags movies
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} models.GetMovieList
// @Router /movies [get]
func (h *MovieHandler) GetList(w http.ResponseWriter, r *http.Request) {
	ctx, span := o11y.Tracer().Start(r.Context(), "MovieHandler.GetList")
	defer span.End()

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	movies, err := h.s.GetList(ctx, page, limit)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, movies)
}

// DeleteMovie godoc
// @Summary Delete a movie
// @Description Delete a movie by ID
// @Tags movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 204 {object} nil
// @Failure 404 {object} types.JsonResponse
// @Router /movies/{id} [delete]
func (h *MovieHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, span := o11y.Tracer().Start(r.Context(), "MovieHandler.Delete")
	defer span.End()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = h.s.Delete(ctx, id)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
