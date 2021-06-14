package event

import (
	"context"
	"fmt"

	"github.com/osapers/mch-back/internal/provider/postgres"
	"github.com/osapers/mch-back/internal/types"
	"github.com/osapers/mch-back/pkg/fakedata"
	"github.com/osapers/mch-back/pkg/pgutil"
)

// Service provides event logic
type Service struct {
	storage *storage
}

func NewService(pgConn *postgres.Conn) *Service {
	return &Service{
		storage: newStorage(pgConn),
	}
}

func (s *Service) List(ctx context.Context) ([]*types.Event, error) {
	events, err := s.storage.list(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}

	return events, nil
}

func (s *Service) LoadTestEvents(count int) error {
	events := fakedata.Events(count)

	_, err := s.storage.conn.DB.Model(&events).Insert()
	err = pgutil.HandleSqlIntegrityViolation(err)

	return err
}
