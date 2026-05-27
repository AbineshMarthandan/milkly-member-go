package entity

import "time"

// Login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

// Register request  
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
	Phone    string `json:"phone" validate:"required"`
}

// Auth response
type AuthResponse struct {
	Token        string          `json:"token"`
	RefreshToken string          `json:"refresh_token,omitempty"`
	User         *MemberResponse `json:"user"`
	ExpiresAt    time.Time       `json:"expires_at"`
}

// Token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
