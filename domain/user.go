package domain

import (
	"time"

	"github.com/wolfsblu/recipe-manager/domain/permissions"
)

type UserDetails struct {
	Email        string
	PasswordHash string
	Locale       string
}

type Permission struct {
	ID   int64
	Name string
	Slug permissions.Slug
}

type Role struct {
	ID          int64
	Name        string
	Permissions []Permission
}

type User struct {
	ID        int64
	Confirmed bool
	Role      Role
	UserDetails
}

type PasswordResetToken struct {
	User      *User
	Token     string
	CreatedAt time.Time
}

type UserRegistration struct {
	User      *User
	Token     string
	CreatedAt time.Time
}
