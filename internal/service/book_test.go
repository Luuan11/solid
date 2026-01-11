package service

import (
	"context"
	"errors"
	"solid/internal/domain"
	"solid/pkg/mocks"
	"testing"
)

func TestBookService_CreateBook(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &mocks.BookRepository{
			CreateFunc: func(ctx context.Context, book *domain.Book) error {
				book.ID = "test-id"
				return nil
			},
		}
		service := NewBookService(repo)

		book, err := service.CreateBook(ctx, "Clean Code", "Robert Martin", "0132350882")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if book.ID != "test-id" {
			t.Errorf("expected ID test-id, got %s", book.ID)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		repo := &mocks.BookRepository{}
		service := NewBookService(repo)

		_, err := service.CreateBook(ctx, "", "Robert Martin", "0132350882")

		if err == nil {
			t.Error("expected validation error")
		}
	})

	t.Run("repository error", func(t *testing.T) {
		repoErr := errors.New("database error")
		repo := &mocks.BookRepository{
			CreateFunc: func(ctx context.Context, book *domain.Book) error {
				return repoErr
			},
		}
		service := NewBookService(repo)

		_, err := service.CreateBook(ctx, "Clean Code", "Robert Martin", "0132350882")

		if err != repoErr {
			t.Errorf("expected repository error, got %v", err)
		}
	})
}

func TestBookService_GetBook(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedBook := &domain.Book{
			ID:     "test-id",
			Title:  "Clean Code",
			Author: "Robert Martin",
			ISBN:   "0132350882",
		}
		repo := &mocks.BookRepository{
			FindByIDFunc: func(ctx context.Context, id string) (*domain.Book, error) {
				return expectedBook, nil
			},
		}
		service := NewBookService(repo)

		book, err := service.GetBook(ctx, "test-id")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if book.ID != expectedBook.ID {
			t.Errorf("expected book ID %s, got %s", expectedBook.ID, book.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := &mocks.BookRepository{
			FindByIDFunc: func(ctx context.Context, id string) (*domain.Book, error) {
				return nil, domain.ErrBookNotFound
			},
		}
		service := NewBookService(repo)

		_, err := service.GetBook(ctx, "nonexistent")

		if err != domain.ErrBookNotFound {
			t.Errorf("expected ErrBookNotFound, got %v", err)
		}
	})
}

func TestBookService_UpdateBook(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		existingBook := &domain.Book{
			ID:     "test-id",
			Title:  "Old Title",
			Author: "Old Author",
			ISBN:   "0132350882",
		}
		repo := &mocks.BookRepository{
			FindByIDFunc: func(ctx context.Context, id string) (*domain.Book, error) {
				return existingBook, nil
			},
			UpdateFunc: func(ctx context.Context, book *domain.Book) error {
				return nil
			},
		}
		service := NewBookService(repo)

		book, err := service.UpdateBook(ctx, "test-id", "New Title", "New Author", "9780134494166")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if book.Title != "New Title" {
			t.Errorf("expected title New Title, got %s", book.Title)
		}
		if book.Author != "New Author" {
			t.Errorf("expected author New Author, got %s", book.Author)
		}
	})

	t.Run("partial update", func(t *testing.T) {
		existingBook := &domain.Book{
			ID:     "test-id",
			Title:  "Old Title",
			Author: "Old Author",
			ISBN:   "0132350882",
		}
		repo := &mocks.BookRepository{
			FindByIDFunc: func(ctx context.Context, id string) (*domain.Book, error) {
				return existingBook, nil
			},
			UpdateFunc: func(ctx context.Context, book *domain.Book) error {
				return nil
			},
		}
		service := NewBookService(repo)

		book, err := service.UpdateBook(ctx, "test-id", "New Title", "", "")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if book.Title != "New Title" {
			t.Errorf("expected title New Title, got %s", book.Title)
		}
		if book.Author != "Old Author" {
			t.Errorf("expected author unchanged, got %s", book.Author)
		}
	})
}

func TestBookService_DeleteBook(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &mocks.BookRepository{
			DeleteFunc: func(ctx context.Context, id string) error {
				return nil
			},
		}
		service := NewBookService(repo)

		err := service.DeleteBook(ctx, "test-id")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := &mocks.BookRepository{
			DeleteFunc: func(ctx context.Context, id string) error {
				return domain.ErrBookNotFound
			},
		}
		service := NewBookService(repo)

		err := service.DeleteBook(ctx, "nonexistent")

		if err != domain.ErrBookNotFound {
			t.Errorf("expected ErrBookNotFound, got %v", err)
		}
	})
}
