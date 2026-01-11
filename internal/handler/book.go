package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"solid/internal/domain"
	"solid/internal/service"
	"time"

	"github.com/gorilla/mux"
)

const requestTimeout = 5 * time.Second

type BookHandler struct {
	service *service.BookService
}

func NewBookHandler(service *service.BookService) *BookHandler {
	return &BookHandler{
		service: service,
	}
}

type createBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
}

type updateBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
}

type errorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

func (h *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	var req createBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload", "INVALID_JSON")
		return
	}

	book, err := h.service.CreateBook(ctx, req.Title, req.Author, req.ISBN)
	if err != nil {
		handleError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, book)
}

func (h *BookHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	id := mux.Vars(r)["id"]

	book, err := h.service.GetBook(ctx, id)
	if err != nil {
		handleError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, book)
}

func (h *BookHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	books, err := h.service.ListBooks(ctx)
	if err != nil {
		handleError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, books)
}

func (h *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	id := mux.Vars(r)["id"]

	var req updateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload", "INVALID_JSON")
		return
	}

	book, err := h.service.UpdateBook(ctx, id, req.Title, req.Author, req.ISBN)
	if err != nil {
		handleError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, book)
}

func (h *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	id := mux.Vars(r)["id"]

	if err := h.service.DeleteBook(ctx, id); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleError(w http.ResponseWriter, err error) {
	statusCode := domain.GetStatusCode(err)
	code := ""
	
	if domainErr, ok := err.(*domain.DomainError); ok {
		code = domainErr.Code
	}
	
	respondWithError(w, statusCode, err.Error(), code)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message, errorCode string) {
	respondWithJSON(w, code, errorResponse{
		Error: message,
		Code:  errorCode,
	})
}
