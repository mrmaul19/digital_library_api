package domain

import (
    "time"
    "github.com/google/uuid"
)

type UserRole string

const (
    RoleAdmin UserRole = "admin"
    RoleUser  UserRole = "user"
)

// Entity
type User struct {
    ID        uuid.UUID `db:"id"        json:"id"`
    Name      string    `db:"name"      json:"name"`
    Email     string    `db:"email"     json:"email"`
    Password  string    `db:"password"  json:"-"`
    Role      UserRole  `db:"role"      json:"role"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// DTOs (Data Transfer Objects)
type RegisterRequest struct {
    Name     string `json:"name"     validate:"required,min=2,max=100"`
    Email    string `json:"email"    validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
    Email    string `json:"email"    validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
    Token string `json:"token"`
    User  *User  `json:"user"`
}

// Repository Interface — usecase tahu "apa" yang bisa dilakukan, bukan "bagaimana"
type UserRepository interface {
    Create(user *User) error
    FindByEmail(email string) (*User, error)
    FindByID(id uuid.UUID) (*User, error)
}

// Usecase Interface
type AuthUsecase interface {
    Register(req *RegisterRequest) (*AuthResponse, error)
    Login(req *LoginRequest) (*AuthResponse, error)
}