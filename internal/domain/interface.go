package domain

import "context"

type BookRepository interface {
	Create(ctx context.Context, book *Book) error
	FindByID(ctx context.Context, id string) (*Book, error)
	FindByISBN(ctx context.Context, isbn string) (*Book, error)
	FindAll(ctx context.Context) ([]*Book, error)
	Update(ctx context.Context, book *Book) error
	Delete(ctx context.Context, id string) error
}
