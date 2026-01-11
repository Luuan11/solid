package domain

import (
	"testing"
)

func TestNewBook(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		author    string
		isbn      string
		wantError bool
	}{
		{
			name:      "valid book with ISBN-10",
			title:     "Clean Code",
			author:    "Robert Martin",
			isbn:      "0132350882",
			wantError: false,
		},
		{
			name:      "valid book with ISBN-13",
			title:     "Clean Architecture",
			author:    "Robert Martin",
			isbn:      "9780134494166",
			wantError: false,
		},
		{
			name:      "valid book with ISBN hyphens",
			title:     "DDD",
			author:    "Eric Evans",
			isbn:      "978-0-321-12521-5",
			wantError: false,
		},
		{
			name:      "empty title",
			title:     "",
			author:    "Robert Martin",
			isbn:      "0132350882",
			wantError: true,
		},
		{
			name:      "empty author",
			title:     "Clean Code",
			author:    "",
			isbn:      "0132350882",
			wantError: true,
		},
		{
			name:      "empty isbn",
			title:     "Clean Code",
			author:    "Robert Martin",
			isbn:      "",
			wantError: true,
		},
		{
			name:      "invalid isbn format",
			title:     "Clean Code",
			author:    "Robert Martin",
			isbn:      "123",
			wantError: true,
		},
		{
			name:      "title too long",
			title:     string(make([]byte, 201)),
			author:    "Robert Martin",
			isbn:      "0132350882",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book, err := NewBook(tt.title, tt.author, tt.isbn)

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if book.Title != tt.title {
				t.Errorf("title = %v, want %v", book.Title, tt.title)
			}
			if book.Author != tt.author {
				t.Errorf("author = %v, want %v", book.Author, tt.author)
			}
			if book.CreatedAt.IsZero() {
				t.Error("CreatedAt should not be zero")
			}
			if book.UpdatedAt.IsZero() {
				t.Error("UpdatedAt should not be zero")
			}
		})
	}
}

func TestValidateISBN(t *testing.T) {
	tests := []struct {
		name      string
		isbn      string
		wantError bool
	}{
		{"valid ISBN-10", "0132350882", false},
		{"valid ISBN-13", "9780134494166", false},
		{"valid with hyphens", "978-0-13-449416-6", false},
		{"empty", "", true},
		{"too short", "123", true},
		{"invalid characters", "012A350882", true},
		{"wrong length", "01323508821", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateISBN(tt.isbn)
			if (err != nil) != tt.wantError {
				t.Errorf("validateISBN() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
