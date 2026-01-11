package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"solid/internal/handler"
	"solid/internal/middleware"
	"solid/internal/repository"
	"solid/internal/service"

	"github.com/gorilla/mux"
)

func main() {
	bookRepository := repository.NewInMemoryBookRepository()
	bookService := service.NewBookService(bookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	router := setupRouter(bookHandler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	gracefulShutdown(srv)
}

func setupRouter(bookHandler *handler.BookHandler) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/health", healthCheck).Methods(http.MethodGet)
	router.HandleFunc("/books", bookHandler.Create).Methods(http.MethodPost)
	router.HandleFunc("/books", bookHandler.List).Methods(http.MethodGet)
	router.HandleFunc("/books/{id}", bookHandler.GetByID).Methods(http.MethodGet)
	router.HandleFunc("/books/{id}", bookHandler.Update).Methods(http.MethodPut)
	router.HandleFunc("/books/{id}", bookHandler.Delete).Methods(http.MethodDelete)

	router.Use(middleware.Recovery)
	router.Use(middleware.Logger)

	return router
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server stopped gracefully")
}
