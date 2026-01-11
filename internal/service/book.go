package service

import (
	"context"
	"solid/internal/domain"
	"time"
)

type BookService struct {
	repository domain.BookRepository
}

func NewBookService(repository domain.BookRepository) *BookService {
	return &BookService{
		repository: repository,
	}
}

func (s *BookService) CreateBook(ctx context.Context, title, author, isbn string) (*domain.Book, error) {
	book, err := domain.NewBook(title, author, isbn)
	if err != nil {
		return nil, err
	}

	if err := s.repository.Create(ctx, book); err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BookService) GetBook(ctx context.Context, id string) (*domain.Book, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *BookService) ListBooks(ctx context.Context) ([]*domain.Book, error) {
	return s.repository.FindAll(ctx)
}

func (s *BookService) UpdateBook(ctx context.Context, id, title, author, isbn string) (*domain.Book, error) {
	book, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if title != "" {
		book.Title = title
	}
	if author != "" {
		book.Author = author
	}
	if isbn != "" {
		book.ISBN = isbn
	}

	book.UpdatedAt = time.Now()

	if err := s.repository.Update(ctx, book); err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BookService) DeleteBook(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}
