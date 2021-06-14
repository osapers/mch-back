package tag

import (
	"context"

	"github.com/go-pg/pg/v10/orm"
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

func (s *storage) search(ctx context.Context, query string, limit int) ([]string, error) {
	var tags []string

	err := s.conn.DB.ModelContext(ctx, (*types.Tag)(nil)).
		Column("name").
		WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			q = q.WhereOr("name ILIKE ?", query).
				WhereOr("name ILIKE ?", query+"%")
			return q, nil
		}).
		Limit(limit).
		Select(&tags)

	if pgutil.IsSqlNoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return tags, nil
}
