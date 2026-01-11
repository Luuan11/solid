package repository

import (
	"context"
	"solid/internal/domain"
	"testing"
)

func TestInMemoryBookRepository_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := NewInMemoryBookRepository()
		book := &domain.Book{
			Title:  "Clean Code",
			Author: "Robert Martin",
			ISBN:   "0132350882",
		}

		err := repo.Create(ctx, book)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if book.ID == "" {
			t.Error("expected ID to be generated")
		}
	})

	t.Run("duplicate ISBN", func(t *testing.T) {
		repo := NewInMemoryBookRepository()
		book1 := &domain.Book{
			Title:  "Clean Code",
			Author: "Robert Martin",
			ISBN:   "0132350882",
		}
		book2 := &domain.Book{
			Title:  "Another Book",
			Author: "Another Author",
			ISBN:   "0132350882",
		}

		repo.Create(ctx, book1)
		err := repo.Create(ctx, book2)

		if err != domain.ErrBookAlreadyExists {
			t.Errorf("expected ErrBookAlreadyExists, got %v", err)
		}
	})
}

func TestInMemoryBookRepository_FindByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := NewInMemoryBookRepository()
		book := &domain.Book{
			Title:  "Clean Code",
			Author: "Robert Martin",
			ISBN:   "0132350882",
		}
		repo.Create(ctx, book)

		found, err := repo.FindByID(ctx, book.ID)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if found.Title != book.Title {
			t.Errorf("expected title %s, got %s", book.Title, found.Title)
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := NewInMemoryBookRepository()

		_, err := repo.FindByID(ctx, "nonexistent")

		if err != domain.ErrBookNotFound {
			t.Errorf("expected ErrBookNotFound, got %v", err)
		}
	})
}

func TestInMemoryBookRepository_FindByISBN(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := NewInMemoryBookRepository()
		book := &domain.Book{
			Title:  "Clean Code",
			Author: "Robert Martin",
			ISBN:   "0132350882",
		}
		repo.Create(ctx, book)

		found, err := repo.FindByISBN(ctx, "0132350882")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if found.ISBN != book.ISBN {
			t.Errorf("expected ISBN %s, got %s", book.ISBN, found.ISBN)
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := NewInMemoryBookRepository()

		_, err := repo.FindByISBN(ctx, "9999999999")

		if err != domain.ErrBookNotFound {
			t.Errorf("expected ErrBookNotFound, got %v", err)
		}
	})
}

func TestInMemoryBookRepository_FindAll(t *testing.T) {
	ctx := context.Background()

	t.Run("empty repository", func(t *testing.T) {
		repo := NewInMemoryBookRepository()

		books, err := repo.FindAll(ctx)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(books) != 0 {
			t.Errorf("expected 0 books, got %d", len(books))
		}
	})

	t.Run("multiple books", func(t *testing.T) {
		repo := NewInMemoryBookRepository()
		book1 := &domain.Book{Title: "Book 1", Author: "Author 1", ISBN: "1111111111"}
		book2 := &domain.Book{Title: "Book 2", Author: "Author 2", ISBN: "2222222222"}
		repo.Create(ctx, book1)
		repo.Create(ctx, book2)

		books, err := repo.FindAll(ctx)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(books) != 2 {
			t.Errorf("expected 2 books, got %d", len(books))
		}
	})
}

func TestInMemoryBookRepository_Update(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := NewInMemoryBookRepository()
		book := &domain.Book{
			Title:  "Clean Code",
			Author: "Robert Martin",
			ISBN:   "0132350882",
		}
		repo.Create(ctx, book)

		book.Title = "Updated Title"
		err := repo.Update(ctx, book)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		found, _ := repo.FindByID(ctx, book.ID)
		if found.Title != "Updated Title" {
			t.Errorf("expected title Updated Title, got %s", found.Title)
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := NewInMemoryBookRepository()
		book := &domain.Book{
			ID:     "nonexistent",
			Title:  "Clean Code",
			Author: "Robert Martin",
			ISBN:   "0132350882",
		}

		err := repo.Update(ctx, book)

		if err != domain.ErrBookNotFound {
			t.Errorf("expected ErrBookNotFound, got %v", err)
		}
	})

	t.Run("duplicate ISBN on update", func(t *testing.T) {
		repo := NewInMemoryBookRepository()
		book1 := &domain.Book{Title: "Book 1", Author: "Author 1", ISBN: "1111111111"}
		book2 := &domain.Book{Title: "Book 2", Author: "Author 2", ISBN: "2222222222"}
		repo.Create(ctx, book1)
		repo.Create(ctx, book2)

		book2.ISBN = "1111111111"
		err := repo.Update(ctx, book2)

		if err != domain.ErrBookAlreadyExists {
			t.Errorf("expected ErrBookAlreadyExists, got %v", err)
		}
	})
}

func TestInMemoryBookRepository_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := NewInMemoryBookRepository()
		book := &domain.Book{
			Title:  "Clean Code",
			Author: "Robert Martin",
			ISBN:   "0132350882",
		}
		repo.Create(ctx, book)

		err := repo.Delete(ctx, book.ID)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		_, err = repo.FindByID(ctx, book.ID)
		if err != domain.ErrBookNotFound {
			t.Error("expected book to be deleted")
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := NewInMemoryBookRepository()

		err := repo.Delete(ctx, "nonexistent")

		if err != domain.ErrBookNotFound {
			t.Errorf("expected ErrBookNotFound, got %v", err)
		}
	})
}
