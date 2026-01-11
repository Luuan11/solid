package mocks

import (
	"context"
	"solid/internal/domain"
)

type BookRepository struct {
	CreateFunc    func(ctx context.Context, book *domain.Book) error
	FindByIDFunc  func(ctx context.Context, id string) (*domain.Book, error)
	FindByISBNFunc func(ctx context.Context, isbn string) (*domain.Book, error)
	FindAllFunc   func(ctx context.Context) ([]*domain.Book, error)
	UpdateFunc    func(ctx context.Context, book *domain.Book) error
	DeleteFunc    func(ctx context.Context, id string) error
}

func (m *BookRepository) Create(ctx context.Context, book *domain.Book) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, book)
	}
	return nil
}

func (m *BookRepository) FindByID(ctx context.Context, id string) (*domain.Book, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, domain.ErrBookNotFound
}

func (m *BookRepository) FindByISBN(ctx context.Context, isbn string) (*domain.Book, error) {
	if m.FindByISBNFunc != nil {
		return m.FindByISBNFunc(ctx, isbn)
	}
	return nil, domain.ErrBookNotFound
}

func (m *BookRepository) FindAll(ctx context.Context) ([]*domain.Book, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}
	return []*domain.Book{}, nil
}

func (m *BookRepository) Update(ctx context.Context, book *domain.Book) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, book)
	}
	return nil
}

func (m *BookRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}
