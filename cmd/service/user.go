package service

import (
	"context"
	"fmt"
	"time"
)

type UserStore interface {
	CreateUser(ctx context.Context, user *User) error
	GetUsers(ctx context.Context) ([]User, error)
	GetUserById(ctx context.Context, id string) (*User, error)
	ActivateUser(ctx context.Context, id string) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUserName(ctx context.Context, id string, firstName, middleName, lastName *string) (*User, error)
}

type UserService interface {
	CreateUser(ctx context.Context, user *User) error
	GetUsers(ctx context.Context) ([]User, error)
	GetUserById(ctx context.Context, id string) (*User, error)
	ActivateUser(ctx context.Context, id string) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUserName(ctx context.Context, id string, firstName, middleName, lastName *string) (*User, error)
}

func NewUserService(
	userStore UserStore,
) UserService {
	return &userService{
		userStore: userStore,
	}
}

type userService struct {
	userStore UserStore
}

type User struct {
	ID         string     `json:"id"`
	FirstName  string     `json:"first_name"`
	MiddleName *string    `json:"middle_name"`
	LastName   string     `json:"last_name"`
	Type       UserType   `json:"type"`
	Status     UserStatus `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type (
	UserType string

	UserStatus string
)

var (
	UserTypeDriver   UserType = "driver"
	UserTypeCustomer UserType = "customer"

	UserTypes = []UserType{
		UserTypeDriver,
		UserTypeCustomer,
	}

	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"

	UserStatuses = []UserStatus{
		UserStatusActive,
		UserStatusInactive,
	}
)

func (v *UserType) UnmarshalText(text []byte) error {
	for _, val := range UserTypes {
		if string(text) == string(val) {
			*v = val
			return nil
		}
	}
	return fmt.Errorf("invalid UserType: %s", string(text))
}

func (v *UserStatus) UnmarshalText(text []byte) error {
	for _, val := range UserStatuses {
		if string(text) == string(val) {
			*v = val
			return nil
		}
	}
	return fmt.Errorf("invalid UserStatus: %s", string(text))
}

func (s *userService) CreateUser(ctx context.Context, user *User) error {
	// suppose to be validation / business logic here
	return s.userStore.CreateUser(ctx, user)
}

func (s *userService) GetUsers(ctx context.Context) ([]User, error) {
	return s.userStore.GetUsers(ctx)
}

func (s *userService) GetUserById(ctx context.Context, id string) (*User, error) {
	user, err := s.userStore.GetUserById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (s *userService) ActivateUser(ctx context.Context, id string) (*User, error) {
	user, err := s.userStore.ActivateUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to activate user: %w", err)
	}
	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	if err := s.userStore.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (s *userService) UpdateUserName(ctx context.Context, id string, firstName, middleName, lastName *string) (*User, error) {
	user, err := s.userStore.UpdateUserName(ctx, id, firstName, middleName, lastName)
	if err != nil {
		return nil, fmt.Errorf("failed to update user name: %w", err)
	}
	return user, nil
}
