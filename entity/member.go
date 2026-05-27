package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Member struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name" validate:"required"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	Phone     string             `bson:"phone" json:"phone"`
	Password  string             `bson:"password" json:"password"`
	Status    MemberStatus       `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type MemberStatus string

const (
	MemberStatusActive    MemberStatus = "ACTIVE"
	MemberStatusInactive  MemberStatus = "INACTIVE"
	MemberStatusSuspended MemberStatus = "SUSPENDED"
	MemberStatusRejected  MemberStatus = "REJECTED"
)

type MemberResponse struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Email     string       `json:"email"`
	Phone     string       `json:"phone"`
	Status    MemberStatus `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func (m *Member) ToResponse() *MemberResponse {
	return &MemberResponse{
		ID:        m.ID.Hex(),
		Name:      m.Name,
		Email:     m.Email,
		Phone:     m.Phone,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}