package tag

import (
	"context"
	"fmt"

	"github.com/osapers/mch-back/internal/provider/postgres"
)

// Service provides tag logic
type Service struct {
	storage *storage
}

func NewService(pgConn *postgres.Conn) *Service {
	return &Service{
		storage: newStorage(pgConn),
	}
}

func (s *Service) Search(ctx context.Context, query string, limit int) ([]string, error) {
	tags, err := s.storage.search(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("search tags: %w", err)
	}

	if len(tags) == 0 {
		return []string{}, nil
	}

	return tags, nil
}
