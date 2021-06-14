package event

import (
	"context"
	"time"

	"github.com/osapers/mch-back/internal/provider/postgres"
	"github.com/osapers/mch-back/internal/types"
	"github.com/osapers/mch-back/pkg/pgutil"
)

// storage is abstraction to manage db data
type storage struct {
	conn *postgres.Conn
}

func newStorage(conn *postgres.Conn) *storage {
	return &storage{
		conn: conn,
	}
}

func (s *storage) list(ctx context.Context, exceptPast bool) ([]*types.Event, error) {
	var events []*types.Event

	query := s.conn.DB.ModelContext(ctx, &events).Order("date asc")

	if exceptPast {
		query = query.Where("date > ?", time.Now().Unix())
	}

	err := query.Select()

	if pgutil.IsSqlNoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return events, nil
}
