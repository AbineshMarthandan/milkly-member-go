package service

import (
	"context"
	"errors"
	"milkly-member/entity"
	"milkly-member/repository"
	"milkly-member/utils"
)

type AuthService interface {
	Register(ctx context.Context, req *entity.RegisterRequest) (*entity.AuthResponse, error)
	Login(ctx context.Context, req *entity.LoginRequest) (*entity.AuthResponse, error)
	RefreshToken(ctx context.Context, req *entity.RefreshTokenRequest) (*entity.AuthResponse, error)
	ValidateToken(token string) (*entity.Member, error)
}

type authService struct {
	memberRepo repository.MemberRepository
}

func NewAuthService(memberRepo repository.MemberRepository) AuthService {
	return &authService{
		memberRepo: memberRepo,
	}
}

func (s *authService) Register(ctx context.Context, req *entity.RegisterRequest) (*entity.AuthResponse, error) {
	member, _ := s.memberRepo.GetByEmail(ctx, req.Email)
	if member != nil {
		return nil, errors.New("MEMBER_ALREADY_EXIST")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("PASSWORD_HASHING_FAILED")
	}
	memberInfo := entity.Member{
		Name: req.Name,
		Email: req.Email,
		Phone: req.Phone,
		Password: hashedPassword,
	}
	createdMemberDetails, err := s.memberRepo.Create(ctx, &memberInfo)
	if err != nil {
		return nil, errors.New("MEMBER_NOT_CREATED_DUE_TO_INTERNAL_ISSUE")
	}
	// TODO: Create member with hashedPassword
	// 4. Generate JWT token
	// 5. Return auth response
	return &entity.AuthResponse{
		Token: "dummy-jwt-token",
		User: createdMemberDetails.ToResponse(),
	}, nil
}

func (s *authService) Login(ctx context.Context, req *entity.LoginRequest) (*entity.AuthResponse, error) {
	// 1. Find member by email
	member, err := s.memberRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("MEMBER_NOT_FOUND")
	}

	// TODO: 2. Verify password (add bcrypt comparison here)
	_ = req.Password // For now, we'll skip password verification
	// TODO: 3. Generate JWT token (implement JWT generation here)
	// 4. Return auth response
	return &entity.AuthResponse{
		Token: "dummy-jwt-token",
		User:  member.ToResponse(), // Use the actual member from database
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req *entity.RefreshTokenRequest) (*entity.AuthResponse, error) {
	// TODO: Implement refresh token logic
	// 1. Validate refresh token
	// 2. Generate new JWT token
	// 3. Return auth response=
	// Placeholder implementation
	return &entity.AuthResponse{
		Token: "new-dummy-jwt-token",
	}, nil
}

func (s *authService) ValidateToken(token string) (*entity.Member, error) {
	// TODO: Implement token validation logic
	// 1. Parse JWT token
	// 2. Validate signature and expiration
	// 3. Extract user ID
	// 4. Return member
	// Placeholder implementation
	return &entity.Member{
		Name:  "Test User",
		Email: "test@example.com",
	}, nil
}
