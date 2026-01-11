package domain

import (
	"regexp"
	"strings"
	"time"
)

const (
	MaxTitleLength  = 200
	MaxAuthorLength = 100
	MaxISBNLength   = 17
)

var isbnRegex = regexp.MustCompile(`^(?:\d{10}|\d{13}|[\d-]{13,17})$`)

type Book struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	ISBN      string    `json:"isbn"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewBook(title, author, isbn string) (*Book, error) {
	if err := validateBook(title, author, isbn); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Book{
		Title:     strings.TrimSpace(title),
		Author:    strings.TrimSpace(author),
		ISBN:      normalizeISBN(isbn),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func validateBook(title, author, isbn string) error {
	if err := validateTitle(title); err != nil {
		return err
	}
	if err := validateAuthor(author); err != nil {
		return err
	}
	if err := validateISBN(isbn); err != nil {
		return err
	}
	return nil
}

func validateTitle(title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return ErrInvalidInput.WithMessage("title cannot be empty")
	}
	if len(title) > MaxTitleLength {
		return ErrInvalidInput.WithMessage("title exceeds maximum length")
	}
	return nil
}

func validateAuthor(author string) error {
	author = strings.TrimSpace(author)
	if author == "" {
		return ErrInvalidInput.WithMessage("author cannot be empty")
	}
	if len(author) > MaxAuthorLength {
		return ErrInvalidInput.WithMessage("author exceeds maximum length")
	}
	return nil
}

func validateISBN(isbn string) error {
	isbn = strings.TrimSpace(isbn)
	if isbn == "" {
		return ErrInvalidInput.WithMessage("isbn cannot be empty")
	}
	if !isbnRegex.MatchString(isbn) {
		return ErrInvalidInput.WithMessage("isbn format is invalid")
	}
	normalized := normalizeISBN(isbn)
	if len(normalized) != 10 && len(normalized) != 13 {
		return ErrInvalidInput.WithMessage("isbn must be 10 or 13 digits")
	}
	return nil
}

func normalizeISBN(isbn string) string {
	return strings.ReplaceAll(strings.TrimSpace(isbn), "-", "")
}
