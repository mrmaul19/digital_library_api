package usecase

import (
    "errors"
    "math"
    "time"

    "github.com/google/uuid"
    "github.com/mrmaul19/digital-library/internal/domain"
)

type bookUsecase struct {
    bookRepo domain.BookRepository
}

func NewBookUsecase(bookRepo domain.BookRepository) domain.BookUsecase {
    return &bookUsecase{bookRepo: bookRepo}
}

func (u *bookUsecase) Create(req *domain.CreateBookRequest, coverURL, fileURL string, uploadedBy uuid.UUID) (*domain.Book, error) {
    now := time.Now()
    book := &domain.Book{
        ID:          uuid.New(),
        Title:       req.Title,
        Author:      req.Author,
        Genre:       req.Genre,
        Description: req.Description,
        CoverURL:    coverURL,
        FileURL:     fileURL,
        Pages:       req.Pages,
        Year:        req.Year,
        UploadedBy:  uploadedBy,
        CreatedAt:   now,
        UpdatedAt:   now,
    }

    if err := u.bookRepo.Create(book); err != nil {
        return nil, errors.New("gagal menyimpan buku")
    }

    return book, nil
}

func (u *bookUsecase) GetByID(id uuid.UUID) (*domain.Book, error) {
    book, err := u.bookRepo.FindByID(id)
    if err != nil {
        return nil, err
    }
    if book == nil {
        return nil, errors.New("buku tidak ditemukan")
    }
    return book, nil
}

func (u *bookUsecase) GetAll(filter domain.BookFilter) (*domain.BookListResponse, error) {
    // Default pagination
    if filter.Page <= 0 {
        filter.Page = 1
    }
    if filter.Limit <= 0 || filter.Limit > 50 {
        filter.Limit = 10
    }

    books, total, err := u.bookRepo.FindAll(filter)
    if err != nil {
        return nil, errors.New("gagal mengambil data buku")
    }

    totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

    return &domain.BookListResponse{
        Books:      books,
        Total:      total,
        Page:       filter.Page,
        Limit:      filter.Limit,
        TotalPages: totalPages,
    }, nil
}

func (u *bookUsecase) Update(id uuid.UUID, req *domain.UpdateBookRequest) (*domain.Book, error) {
    book, err := u.bookRepo.FindByID(id)
    if err != nil {
        return nil, err
    }
    if book == nil {
        return nil, errors.New("buku tidak ditemukan")
    }

    // Hanya update field yang dikirim (pointer = opsional)
    if req.Title != nil {
        book.Title = *req.Title
    }
    if req.Author != nil {
        book.Author = *req.Author
    }
    if req.Genre != nil {
        book.Genre = *req.Genre
    }
    if req.Description != nil {
        book.Description = *req.Description
    }
    if req.Pages != nil {
        book.Pages = *req.Pages
    }
    if req.Year != nil {
        book.Year = *req.Year
    }

    book.UpdatedAt = time.Now()

    if err := u.bookRepo.Update(book); err != nil {
        return nil, errors.New("gagal mengupdate buku")
    }

    return book, nil
}

func (u *bookUsecase) Delete(id uuid.UUID) error {
    book, err := u.bookRepo.FindByID(id)
    if err != nil {
        return err
    }
    if book == nil {
        return errors.New("buku tidak ditemukan")
    }

    return u.bookRepo.Delete(id)
}