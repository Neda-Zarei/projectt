package dto

import "time"

type LoginRequest struct {
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required"`
	CaptchaToken string `json:"captcha_token" validate:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type ListUsersResponse struct {
	Users      []UserResponse `json:"users"`
	Pagination Pagination     `json:"pagination"`
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

// Plan represents a subscription plan
type PlanResponse struct {
	ID        string `json:"id" example:"plan_001"`
	Name      string `json:"name" example:"Premium"`
	StartDate string `json:"startDate" example:"2025-01-01"`
	EndDate   string `json:"endDate" example:"2025-12-31"`
}

// Error response
type Error struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"invalid request"`
}

type ToggleUserActiveRequest struct {
    Active bool `json:"active"`
}