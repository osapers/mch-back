package postgres

import (
	"context"

	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
)

type debugHook struct {
	logger  *zap.Logger
	verbose bool
}

var _ pg.QueryHook = (*debugHook)(nil)

func (h debugHook) BeforeQuery(ctx context.Context, evt *pg.QueryEvent) (context.Context, error) {
	q, err := evt.FormattedQuery()
	if err != nil {
		return nil, err
	}

	if evt.Err != nil {
		h.logger.Error("", zap.String("query", string(q)), zap.Error(evt.Err))
	} else if h.verbose {
		h.logger.Info("", zap.String("query", string(q)))
	}

	return ctx, nil
}

func (debugHook) AfterQuery(context.Context, *pg.QueryEvent) error {
	return nil
}
