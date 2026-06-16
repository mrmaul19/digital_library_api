package jwtpkg

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "github.com/mrmaul19/digital-library/internal/domain"
)

type Claims struct {
    UserID uuid.UUID       `json:"user_id"`
    Email  string          `json:"email"`
    Role   domain.UserRole `json:"role"`
    jwt.RegisteredClaims
}

type JWTService struct {
    secret      []byte
    expireHours int
}

func NewJWTService(secret string, expireHours int) *JWTService {
    return &JWTService{
        secret:      []byte(secret),
        expireHours: expireHours,
    }
}

func (j *JWTService) GenerateToken(user *domain.User) (string, error) {
    claims := &Claims{
        UserID: user.ID,
        Email:  user.Email,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expireHours) * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.secret)
}

func (j *JWTService) ValidateToken(tokenStr string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return j.secret, nil
    })
    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }

    return claims, nil
}