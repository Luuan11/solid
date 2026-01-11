package repository

import (
	"context"
	"sync"

	"solid/internal/domain"

	"github.com/google/uuid"
)

type InMemoryBookRepository struct {
	mu    sync.RWMutex
	books map[string]*domain.Book
	isbn  map[string]string
}

func NewInMemoryBookRepository() *InMemoryBookRepository {
	return &InMemoryBookRepository{
		books: make(map[string]*domain.Book),
		isbn:  make(map[string]string),
	}
}

func (r *InMemoryBookRepository) Create(ctx context.Context, book *domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.isbn[book.ISBN]; exists {
		return domain.ErrBookAlreadyExists
	}

	book.ID = uuid.New().String()

	bookCopy := *book
	r.books[book.ID] = &bookCopy
	r.isbn[book.ISBN] = book.ID
	return nil
}

func (r *InMemoryBookRepository) FindByID(ctx context.Context, id string) (*domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	book, exists := r.books[id]
	if !exists {
		return nil, domain.ErrBookNotFound
	}
	return book, nil
}

func (r *InMemoryBookRepository) FindByISBN(ctx context.Context, isbn string) (*domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exists := r.isbn[isbn]
	if !exists {
		return nil, domain.ErrBookNotFound
	}
	return r.books[id], nil
}

func (r *InMemoryBookRepository) FindAll(ctx context.Context) ([]*domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	books := make([]*domain.Book, 0, len(r.books))
	for _, book := range r.books {
		books = append(books, book)
	}
	return books, nil
}

func (r *InMemoryBookRepository) Update(ctx context.Context, book *domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.books[book.ID]
	if !exists {
		return domain.ErrBookNotFound
	}

	if existing.ISBN != book.ISBN {
		if _, isbnExists := r.isbn[book.ISBN]; isbnExists {
			return domain.ErrBookAlreadyExists
		}
		delete(r.isbn, existing.ISBN)
		r.isbn[book.ISBN] = book.ID
	}

	bookCopy := *book
	r.books[book.ID] = &bookCopy
	return nil
}

func (r *InMemoryBookRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	book, exists := r.books[id]
	if !exists {
		return domain.ErrBookNotFound
	}

	delete(r.isbn, book.ISBN)
	delete(r.books, id)
	return nil
}
