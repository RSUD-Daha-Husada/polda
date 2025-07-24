package dto

import "github.com/google/uuid"

type CreateUserInput struct {
	Name           string      `json:"name" validate:"required"`
	Username       string      `json:"username" validate:"required"`
	Password       string      `json:"password" validate:"required"`
	Gender         string      `json:"gender"`
	Email          string      `json:"email"`
	Telephone      string      `json:"telephone"`
	Photo          string      `json:"photo"`
	RoleID         uuid.UUID   `json:"role_id" validate:"required"`
	IsActive       bool        `json:"is_active"`
	CreatedBy      uuid.UUID   `json:"created_by"`
	ApplicationIDs []uuid.UUID `json:"application_ids"`
}

type EditUserInput struct {
	UserID         uuid.UUID   `json:"user_id" validate:"required"`
	Name           string      `json:"name" validate:"required"`
	Username       string      `json:"username" validate:"required"`
	Password       *string     `json:"password,omitempty"` // pointer agar bisa bedakan antara kosong dan tidak dikirim
	Gender         string      `json:"gender"`
	Email          string      `json:"email"`
	Telephone      string      `json:"telephone"`
	Photo          string      `json:"photo"`
	RoleID         uuid.UUID   `json:"role_id" validate:"required"`
	IsActive       bool        `json:"is_active"`
	ApplicationIDs []uuid.UUID `json:"application_ids"`
	UpdatedBy      uuid.UUID   `json:"updated_by"`
}

type UpdateProfileInput struct {
	UserID      uuid.UUID
	OldPassword string
	Password    string
	Photo        string
	Telephone   string
}