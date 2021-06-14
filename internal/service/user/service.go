package user

import (
	"context"
	"fmt"

	"github.com/osapers/mch-back/internal/constant"
	"github.com/osapers/mch-back/internal/provider/postgres"
	"github.com/osapers/mch-back/internal/types"
)

// Service provides event logic
type Service struct {
	storage *storage
	kwe     KeywordExtractor
}

type KeywordExtractor interface {
	Extract(ctx context.Context, text string) ([]string, error)
}

// NewService returns new user Service
func NewService(pgConn *postgres.Conn, kwe KeywordExtractor) *Service {
	return &Service{
		storage: newStorage(pgConn),
		kwe:     kwe,
	}
}

// GetByID returns user by ID
func (s *Service) GetByID(ctx context.Context, id string) (*types.User, error) {
	user, err := s.storage.getByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

// Authorize checks that user exists and provided password is correct
// User will be created if doesn't exist
func (s *Service) Authorize(ctx context.Context, email, password string) (*types.User, error) {
	user, err := s.storage.getByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	if user != nil && user.Password != password {
		return nil, newWrongPasswordError()
	}

	if user != nil {
		return user, nil
	}

	user, err = s.storage.create(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

// GetParticipatedEvents returns list of eventID where user is participant
func (s *Service) GetParticipatedEvents(ctx context.Context) (map[string]struct{}, error) {
	userID := constant.GetUserIDFromCtx(ctx)

	events, err := s.storage.getParticipatedEvents(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get events with user participation: %w", err)
	}

	return events, nil
}

func (s *Service) ParticipateInEvent(ctx context.Context, eventID string) error {
	userID := constant.GetUserIDFromCtx(ctx)

	err := s.storage.participateInEvent(ctx, userID, eventID)
	if err != nil {
		return fmt.Errorf("participate in event: %w", err)
	}

	return nil
}

func (s *Service) GetMe(ctx context.Context) (*types.User, error) {
	userID := constant.GetUserIDFromCtx(ctx)

	user, err := s.storage.getByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get me: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user doesn't exists")
	}

	return user, nil
}

func (s *Service) Update(ctx context.Context, reqUser *types.User) (*types.User, error) {
	reqUser.ID = constant.GetUserIDFromCtx(ctx)

	user, err := s.storage.update(ctx, reqUser)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}
