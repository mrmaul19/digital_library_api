package repository

import (
    "database/sql"
    "errors"
    "fmt"
    "strings"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    "github.com/mrmaul19/digital-library/internal/domain"
)

type bookRepository struct {
    db *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) domain.BookRepository {
    return &bookRepository{db: db}
}

func (r *bookRepository) Create(book *domain.Book) error {
    query := `
        INSERT INTO books (id, title, author, genre, description, cover_url, file_url, pages, year, uploaded_by, created_at, updated_at)
        VALUES (:id, :title, :author, :genre, :description, :cover_url, :file_url, :pages, :year, :uploaded_by, :created_at, :updated_at)
    `
    _, err := r.db.NamedExec(query, book)
    return err
}

func (r *bookRepository) FindByID(id uuid.UUID) (*domain.Book, error) {
    var book domain.Book
    query := `SELECT * FROM books WHERE id = $1 LIMIT 1`

    err := r.db.Get(&book, query, id)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &book, nil
}

func (r *bookRepository) FindAll(filter domain.BookFilter) ([]*domain.Book, int, error) {
    // Build query dinamis berdasarkan filter
    conditions := []string{"1=1"}
    args := []interface{}{}
    argIdx := 1

    if filter.Search != "" {
        conditions = append(conditions, fmt.Sprintf(
            "(title ILIKE $%d OR author ILIKE $%d)", argIdx, argIdx+1,
        ))
        searchTerm := "%" + filter.Search + "%"
        args = append(args, searchTerm, searchTerm)
        argIdx += 2
    }

    if filter.Genre != "" {
        conditions = append(conditions, fmt.Sprintf("genre ILIKE $%d", argIdx))
        args = append(args, "%"+filter.Genre+"%")
        argIdx++
    }

    if filter.Author != "" {
        conditions = append(conditions, fmt.Sprintf("author ILIKE $%d", argIdx))
        args = append(args, "%"+filter.Author+"%")
        argIdx++
    }

    if filter.Year > 0 {
        conditions = append(conditions, fmt.Sprintf("year = $%d", argIdx))
        args = append(args, filter.Year)
        argIdx++
    }

    where := "WHERE " + strings.Join(conditions, " AND ")

    // Hitung total dulu (untuk pagination)
    var total int
    countQuery := fmt.Sprintf("SELECT COUNT(*) FROM books %s", where)
    if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
        return nil, 0, err
    }

    // Pagination
    offset := (filter.Page - 1) * filter.Limit
    dataQuery := fmt.Sprintf(
        "SELECT * FROM books %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
        where, argIdx, argIdx+1,
    )
    args = append(args, filter.Limit, offset)

    var books []*domain.Book
    if err := r.db.Select(&books, dataQuery, args...); err != nil {
        return nil, 0, err
    }

    return books, total, nil
}

func (r *bookRepository) Update(book *domain.Book) error {
    query := `
        UPDATE books SET
            title       = :title,
            author      = :author,
            genre       = :genre,
            description = :description,
            cover_url   = :cover_url,
            file_url    = :file_url,
            pages       = :pages,
            year        = :year,
            updated_at  = :updated_at
        WHERE id = :id
    `
    _, err := r.db.NamedExec(query, book)
    return err
}

func (r *bookRepository) Delete(id uuid.UUID) error {
    _, err := r.db.Exec("DELETE FROM books WHERE id = $1", id)
    return err
}