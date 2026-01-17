# 📚 Book Management API

Welcome to the Book Management API! A robust and scalable RESTful API built with Go, following SOLID principles and clean architecture patterns.

<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/REST_API-FF6C37?style=for-the-badge&logo=postman&logoColor=white" alt="REST API" />
  <img src="https://img.shields.io/badge/Clean_Architecture-8B5CF6?style=for-the-badge" alt="Clean Architecture" />
</p>

## ✨ About

This API was designed to demonstrate modern Go development practices with a focus on maintainability, testability, and scalability. The project showcases:

- **Clean Architecture** - Separation of concerns with domain-driven design
- **SOLID Principles** - Following best practices for object-oriented design
- **Comprehensive Testing** - Unit tests with mocks for reliable code
- **Error Handling** - Custom error types for better error management
- **Middleware Support** - Logging and recovery middleware
- **Type Safety** - Leveraging Go's strong typing system

## Architecture

```
cmd/
  api/
    main.go                 # Entry point and DI container
internal/
  domain/
    book.go                 # Domain entity
    repository.go           # Repository interface
  handler/
    book_handler.go         # HTTP handlers
  middleware/
    logger.go               # Logging middleware
    recovery.go             # Panic recovery middleware
  repository/
    book_repository.go      # In-memory implementation
  service/
    book_service.go         # Business rules
```

## SOLID Principles Applied

### S - Single Responsibility Principle

Each component has a single well-defined responsibility:


```go
// book_service.go - ONLY business logic
func (s *BookService) CreateBook(title, author, isbn string) (*domain.Book, error) {
    book, err := domain.NewBook(title, author, isbn)
    if err != nil {
        return nil, err
    }
    return book, s.repository.Create(book)
}
```

### O - Open/Closed Principle

Code is open for extension but closed for modification:

```go
// Easy to add new implementations
type PostgreSQLBookRepository struct { ... }
type MongoBookRepository struct { ... }

// All implement the same interface
type BookRepository interface {
    Create(book *Book) error
    FindByID(id string) (*Book, error)
    // ...
}
```

### L - Liskov Substitution Principle

Any `BookRepository` implementation can replace another without breaking the code:

```go
// Both work perfectly in BookService
bookService := service.NewBookService(repository.NewInMemoryBookRepository())
bookService := service.NewBookService(repository.NewPostgreSQLRepository())
```

### I - Interface Segregation Principle

Lean interfaces with only necessary methods:

```go
// Specific interface for book operations
type BookRepository interface {
    Create(book *Book) error
    FindByID(id string) (*Book, error)
    FindAll() ([]*Book, error)
    Update(book *Book) error
    Delete(id string) error
}
```

### D - Dependency Inversion Principle

High-level modules don't depend on low-level modules, both depend on abstractions:

```go
// Service depends on the INTERFACE, not the concrete implementation
type BookService struct {
    repository domain.BookRepository  // Interface, not concrete struct
}

// Manual dependency injection in main.go
bookRepository := repository.NewInMemoryBookRepository()
bookService := service.NewBookService(bookRepository)
bookHandler := handler.NewBookHandler(bookService)
```

## API Endpoints

### Create Book
```bash
POST /books
Content-Type: application/json

{
  "title": "Clean Code",
  "author": "Robert C. Martin",
  "isbn": "978-0132350884"
}
```

### List All Books
```bash
GET /books
```

### Get Book by ID
```bash
GET /books/{id}
```

### Update Book
```bash
PUT /books/{id}
Content-Type: application/json

{
  "title": "Clean Architecture",
  "author": "Robert C. Martin",
  "isbn": "978-0134494166"
}
```

### Delete Book
```bash
DELETE /books/{id}
```

The API will be available at `http://localhost:8080`

## Usage Examples

```bash
# Create a book
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Domain-Driven Design",
    "author": "Eric Evans",
    "isbn": "978-0321125215"
  }'

# List all books
curl http://localhost:8080/books

# Get specific book
curl http://localhost:8080/books/{id}

# Update book
curl -X PUT http://localhost:8080/books/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "title": "DDD Distilled",
    "author": "Vaughn Vernon",
    "isbn": "978-0134434421"
  }'

# Delete book
curl -X DELETE http://localhost:8080/books/{id}
```

## Response Structure

### Success
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Clean Code",
  "author": "Robert C. Martin",
  "isbn": "978-0132350884",
  "created_at": "2026-01-10T10:30:00Z",
  "updated_at": "2026-01-10T10:30:00Z"
}
```

### Error
```json
{
  "error": "book not found"
}
```

## 📦 Installation

```bash
# Clone repository 
git clone https://github.com/Luuan11/solid.git 

# Install dependencies
go mod tidy

# Run application
go run cmd/api/main.go
```
---
Made with 💜 by [Luan Fernando](https://www.linkedin.com/in/luan-fernando/).