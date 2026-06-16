package domain

import (
    "time"
    "github.com/google/uuid"
)

type Book struct {
    ID          uuid.UUID  `db:"id"           json:"id"`
    Title       string     `db:"title"        json:"title"`
    Author      string     `db:"author"       json:"author"`
    Genre       string     `db:"genre"        json:"genre"`
    Description string     `db:"description"  json:"description"`
    CoverURL    string     `db:"cover_url"    json:"cover_url"`
    FileURL     string     `db:"file_url"     json:"file_url"`
    Pages       int        `db:"pages"        json:"pages"`
    Year        int        `db:"year"         json:"year"`
    UploadedBy  uuid.UUID  `db:"uploaded_by"  json:"uploaded_by"`
    CreatedAt   time.Time  `db:"created_at"   json:"created_at"`
    UpdatedAt   time.Time  `db:"updated_at"   json:"updated_at"`
}

// Filter untuk search
type BookFilter struct {
    Search string
    Genre  string
    Author string
    Year   int
    Page   int
    Limit  int
}

type BookListResponse struct {
    Books      []*Book `json:"books"`
    Total      int     `json:"total"`
    Page       int     `json:"page"`
    Limit      int     `json:"limit"`
    TotalPages int     `json:"total_pages"`
}

type CreateBookRequest struct {
    Title       string `form:"title"       validate:"required,min=1,max=255"`
    Author      string `form:"author"      validate:"required,min=1,max=255"`
    Genre       string `form:"genre"       validate:"required"`
    Description string `form:"description"`
    Pages       int    `form:"pages"       validate:"min=0"`
    Year        int    `form:"year"        validate:"min=1000,max=2100"`
}

type UpdateBookRequest struct {
    Title       *string `json:"title"       validate:"omitempty,min=1,max=255"`
    Author      *string `json:"author"      validate:"omitempty,min=1,max=255"`
    Genre       *string `json:"genre"`
    Description *string `json:"description"`
    Pages       *int    `json:"pages"`
    Year        *int    `json:"year"`
}

// Repository Interface
type BookRepository interface {
    Create(book *Book) error
    FindByID(id uuid.UUID) (*Book, error)
    FindAll(filter BookFilter) ([]*Book, int, error)
    Update(book *Book) error
    Delete(id uuid.UUID) error
}

// Usecase Interface
type BookUsecase interface {
    Create(req *CreateBookRequest, coverURL, fileURL string, uploadedBy uuid.UUID) (*Book, error)
    GetByID(id uuid.UUID) (*Book, error)
    GetAll(filter BookFilter) (*BookListResponse, error)
    Update(id uuid.UUID, req *UpdateBookRequest) (*Book, error)
    Delete(id uuid.UUID) error
}