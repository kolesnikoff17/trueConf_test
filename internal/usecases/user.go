package usecases

import (
	"context"
	"errors"
	"fmt"
	"users/internal/entity"
)

// UserRepo is an interface for repo layer
type UserRepo interface {
	GetByID(ctx context.Context, id int) (entity.User, error)
	Create(ctx context.Context, user entity.User) (int, error)
	Update(ctx context.Context, user entity.User) error
	Delete(ctx context.Context, id int) error
}

// UserUsecase implements http.UserUsecase
type UserUsecase struct {
	r UserRepo
}

// New is a constructor for UserUsecase
func New(r UserRepo) *UserUsecase {
	return &UserUsecase{
		r: r,
	}
}

// GetUserByID return user by its id, entity.ErrNoID if there is no one
func (u *UserUsecase) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	user, err := u.r.GetByID(ctx, id)
	switch {
	case errors.Is(err, entity.ErrNoID):
		return entity.User{}, err
	case err != nil:
		return entity.User{}, fmt.Errorf("UserUsecase - GetUserByID: %w", err)
	}
	return user, nil
}

// CreateUser adds new user, return its id
func (u *UserUsecase) CreateUser(ctx context.Context, user entity.User) (int, error) {
	id, err := u.r.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("UserUsecase - CreateUser: %w", err)
	}
	return id, nil
}

// UpdateUser update user info, return entity.ErrNoID if there is no one
func (u *UserUsecase) UpdateUser(ctx context.Context, user entity.User) error {
	_, err := u.GetUserByID(ctx, user.ID)
	switch {
	case errors.Is(err, entity.ErrNoID):
		return err
	case err != nil:
		return fmt.Errorf("UserUsecase - UpdateUser: %w", err)
	}
	err = u.r.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("UserUsecase - UpdateUser: %w", err)
	}
	return nil
}

// DeleteUser delete user
func (u *UserUsecase) DeleteUser(ctx context.Context, id int) error {
	_, err := u.GetUserByID(ctx, id)
	switch {
	case errors.Is(err, entity.ErrNoID):
		return err
	case err != nil:
		return fmt.Errorf("UserUsecase - DeleteUser: %w", err)
	}
	err = u.r.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("UserUsecase - DeleteUser: %w", err)
	}
	return nil
}
