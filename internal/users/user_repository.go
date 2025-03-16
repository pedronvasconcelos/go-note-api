package users

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Add(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func (u *userRepository) Add(ctx context.Context, user *User) error {
	return u.db.WithContext(ctx).Create(user).Error
}

func (u *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := u.db.WithContext(ctx).Model(&User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	err := u.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) Update(ctx context.Context, user *User) error {
	return u.db.WithContext(ctx).Save(user).Error
}

func NewUserRepository(db *gorm.DB) Repository {
	return &userRepository{db: db}
}
