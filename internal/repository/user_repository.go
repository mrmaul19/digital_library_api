package repository

import (
    "database/sql"
    "errors"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    "github.com/mrmaul19/digital-library/internal/domain"
)

type userRepository struct {
    db *sqlx.DB
}

// Constructor — return interface, bukan struct (clean architecture)
func NewUserRepository(db *sqlx.DB) domain.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
    query := `
        INSERT INTO users (id, name, email, password, role, created_at, updated_at)
        VALUES (:id, :name, :email, :password, :role, :created_at, :updated_at)
    `
    _, err := r.db.NamedExec(query, user)
    return err
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
    var user domain.User
    query := `SELECT * FROM users WHERE email = $1 LIMIT 1`

    err := r.db.Get(&user, query, email)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil // tidak ditemukan, bukan error
    }
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) FindByID(id uuid.UUID) (*domain.User, error) {
    var user domain.User
    query := `SELECT * FROM users WHERE id = $1 LIMIT 1`

    err := r.db.Get(&user, query, id)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &user, nil
}