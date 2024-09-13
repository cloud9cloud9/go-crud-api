package user

import (
	"context"
	"go-rest-api/internal/apperror"
	"go-rest-api/pkg/logging"
)

type Service struct {
	storage Storage
	logger  *logging.Logger
}

func NewService(storage Storage, logger *logging.Logger) *Service {
	return &Service{
		storage: storage,
		logger:  logger,
	}
}

func (s *Service) Create(ctx context.Context, userDTO CreateUserDTO) (User, error) {
	user := User{
		Username:     userDTO.Username,
		Email:        userDTO.Email,
		PasswordHash: userDTO.Password,
	}
	_, err := s.storage.Create(ctx, user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *Service) FindOne(ctx context.Context, id string) (User, error) {
	user, err := s.storage.FindOne(ctx, id)
	if err != nil {
		return User{}, apperror.ErrNotFound
	}
	return user, nil
}

func (s *Service) FindAll(ctx context.Context) ([]User, error) {
	users, err := s.storage.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) Update(ctx context.Context, userDTO UpdateUserDTO) error {
	user := User{
		ID:           userDTO.ID,
		Username:     userDTO.Username,
		Email:        userDTO.Email,
		PasswordHash: userDTO.Password,
	}
	err := s.storage.Update(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	err := s.storage.Delete(ctx, id)
	if err != nil {
		return apperror.ErrNotFound
	}
	return nil
}
