package users

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	FirstName    string
	LastName     string
	Email        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	PremiumUntil *time.Time
	isActive     bool
}

func newUser(firstName, lastName, email string) *User {
	return &User{
		ID:           uuid.New(),
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		PremiumUntil: nil,
		isActive:     true,
	}
}

func (u *User) isPremium() bool {
	return u.PremiumUntil != nil && u.PremiumUntil.After(time.Now())
}

func (u *User) activatePremium30Days() {
	now := time.Now()
	premiumUntil := now.Add(30 * 24 * time.Hour)
	u.PremiumUntil = &premiumUntil
	u.UpdatedAt = now
}
