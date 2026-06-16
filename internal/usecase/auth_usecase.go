package usecase

import (
    "errors"
    "time"

    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "github.com/mrmaul19/digital-library/internal/domain"
    jwtpkg "github.com/mrmaul19/digital-library/pkg/jwt"
)

type authUsecase struct {
    userRepo   domain.UserRepository
    jwtService *jwtpkg.JWTService
}

func NewAuthUsecase(userRepo domain.UserRepository, jwtService *jwtpkg.JWTService) domain.AuthUsecase {
    return &authUsecase{
        userRepo:   userRepo,
        jwtService: jwtService,
    }
}

func (u *authUsecase) Register(req *domain.RegisterRequest) (*domain.AuthResponse, error) {
    // Cek apakah email sudah terdaftar
    existing, err := u.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, err
    }
    if existing != nil {
        return nil, errors.New("email sudah terdaftar")
    }

    // Hash password — JANGAN simpan plain text
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, errors.New("gagal memproses password")
    }

    now := time.Now()
    user := &domain.User{
        ID:        uuid.New(),
        Name:      req.Name,
        Email:     req.Email,
        Password:  string(hashedPassword),
        Role:      domain.RoleUser,
        CreatedAt: now,
        UpdatedAt: now,
    }

    if err := u.userRepo.Create(user); err != nil {
        return nil, errors.New("gagal membuat akun")
    }

    token, err := u.jwtService.GenerateToken(user)
    if err != nil {
        return nil, errors.New("gagal membuat token")
    }

    return &domain.AuthResponse{Token: token, User: user}, nil
}

func (u *authUsecase) Login(req *domain.LoginRequest) (*domain.AuthResponse, error) {
    user, err := u.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, err
    }

    // Pesan error dibuat sama agar tidak bocorkan info "email tidak ada"
    if user == nil {
        return nil, errors.New("email atau password salah")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return nil, errors.New("email atau password salah")
    }

    token, err := u.jwtService.GenerateToken(user)
    if err != nil {
        return nil, errors.New("gagal membuat token")
    }

    return &domain.AuthResponse{Token: token, User: user}, nil
}