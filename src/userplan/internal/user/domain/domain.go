package domain

import (
	"time"
)

type Basic struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"<-:create;"`
	UpdatedAt time.Time
}

type User struct {
	Basic
	OauthID                string `json:"oauth_id,omitempty"`
	Email                  string `json:"email" gorm:"not null"`
	Name                   string `json:"name" gorm:"not null"`
	Password               string `json:"password,omitempty"`
	Phone                  string `json:"phone,omitempty"`
	CompanyName            string `json:"company_name,omitempty"`
	JobTitle               string `json:"job_title,omitempty"`
	Active                 bool   `json:"active" gorm:"not null;default:false"`
	SubscribeNews          bool   `json:"subscribe_news" gorm:"not null;default:true"`
	SubscribeNotifications bool   `json:"subscribe_notifications" gorm:"not null;default:true"`
}

type UserFilter struct {
	UserID uint
	Name   string
	Email  string
	Phone  string
}

type PaginatedUsers struct {
	Users             []*User
	Page, Size, Total int64
}

type UserActivation struct {
	UserID uint
	Active bool
}
