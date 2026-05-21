package service

import (
	"context"
	"errors"
	"fmt"

	"milkly-member/entity"
	"milkly-member/repository"

	"go.mongodb.org/mongo-driver/mongo"
)

type MemberService interface {
	CreateMember(ctx context.Context, req *entity.CreateMemberRequest) (*entity.MemberResponse, error)
	GetMember(ctx context.Context, id string) (*entity.MemberResponse, error)
	UpdateMember(ctx context.Context, id string, req *entity.UpdateMemberRequest) (*entity.MemberResponse, error)
	DeleteMember(ctx context.Context, id string) error
	ListMembers(ctx context.Context, limit, offset int) ([]*entity.MemberResponse, error)
}

type memberService struct {
	memberRepo repository.MemberRepository
}

func NewMemberService(memberRepo repository.MemberRepository) MemberService {
	return &memberService{
		memberRepo: memberRepo,
	}
}

func (s *memberService) CreateMember(ctx context.Context, req *entity.CreateMemberRequest) (*entity.MemberResponse, error) {
	// Check if email already exists
	_, err := s.memberRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}
	if err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("error checking email: %w", err)
	}

	member := &entity.Member{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	createdMember, err := s.memberRepo.Create(ctx, member)
	if err != nil {
		return nil, fmt.Errorf("error creating member: %w", err)
	}

	return createdMember.ToResponse(), nil
}

func (s *memberService) GetMember(ctx context.Context, id string) (*entity.MemberResponse, error) {
	member, err := s.memberRepo.GetByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("member not found")
		}
		return nil, fmt.Errorf("error getting member: %w", err)
	}

	return member.ToResponse(), nil
}

func (s *memberService) UpdateMember(ctx context.Context, id string, req *entity.UpdateMemberRequest) (*entity.MemberResponse, error) {
	// Get existing member
	existingMember, err := s.memberRepo.GetByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("member not found")
		}
		return nil, fmt.Errorf("error getting member: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		existingMember.Name = *req.Name
	}
	if req.Email != nil {
		// Check if new email already exists (but not for the same member)
		memberWithEmail, err := s.memberRepo.GetByEmail(ctx, *req.Email)
		if err == nil && memberWithEmail.ID.Hex() != id {
			return nil, errors.New("email already exists")
		}
		if err != nil && err != mongo.ErrNoDocuments {
			return nil, fmt.Errorf("error checking email: %w", err)
		}
		existingMember.Email = *req.Email
	}
	if req.Phone != nil {
		existingMember.Phone = *req.Phone
	}
	if req.Status != nil {
		existingMember.Status = *req.Status
	}

	updatedMember, err := s.memberRepo.Update(ctx, id, existingMember)
	if err != nil {
		return nil, fmt.Errorf("error updating member: %w", err)
	}

	return updatedMember.ToResponse(), nil
}

func (s *memberService) DeleteMember(ctx context.Context, id string) error {
	// Check if member exists
	_, err := s.memberRepo.GetByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("member not found")
		}
		return fmt.Errorf("error getting member: %w", err)
	}

	err = s.memberRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting member: %w", err)
	}

	return nil
}

func (s *memberService) ListMembers(ctx context.Context, limit, offset int) ([]*entity.MemberResponse, error) {
	members, err := s.memberRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing members: %w", err)
	}

	var responses []*entity.MemberResponse
	for _, member := range members {
		responses = append(responses, member.ToResponse())
	}

	return responses, nil
}