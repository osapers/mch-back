package user

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/osapers/mch-back/internal/types"
	"github.com/osapers/mch-back/pkg/pgutil"
)

func (s *storage) createProject(ctx context.Context, p *types.Project) error {
	_, err := s.conn.DB.ModelContext(ctx, p).Returning("*").Insert()
	err = pgutil.HandleSqlIntegrityViolation(err)

	return err
}

func (s *storage) getRandomTags(ctx context.Context) ([]string, error) { // nolint:unused
	var tags []*types.Tag

	err := s.conn.DB.ModelContext(ctx, &tags).
		OrderExpr("random()").
		Limit(10).
		Select(&tags)

	if pgutil.IsSqlNoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	res := make([]string, 10)

	for i, tag := range tags {
		res[i] = tag.Name
	}

	return res, nil
}

func (s *storage) searchProjects(ctx context.Context, userID string, viewedProjects []string) ([]*types.Project, error) {
	var projects []*types.Project

	q := s.conn.DB.ModelContext(ctx, &projects).Where("owner_id != ?", userID)

	if len(viewedProjects) != 0 {
		q = q.Where("id not in (?)", pg.In(viewedProjects))
	}

	err := q.Select()

	if pgutil.IsSqlNoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *storage) getProjectTags(ctx context.Context, ownerID, projectID string) ([]string, error) {
	var tags []string

	err := s.conn.DB.ModelContext(ctx, (*types.Project)(nil)).
		Column("tags").
		Where("id = ?", projectID).
		Where("owner_id = ?", ownerID).
		Select(&tags)

	if pgutil.IsSqlNoRows(err) {
		return nil, fmt.Errorf("no result")
	}

	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (s *storage) getProjectViewers(ctx context.Context, projectID string) ([]string, error) {
	var users []string

	err := s.conn.DB.ModelContext(ctx, (*types.UserToProject)(nil)).
		Column("user_id").
		Where("project_id = ?", projectID).
		WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			q = q.WhereOr("applied != ?", true).
				WhereOr("confirmed != ?", true).
				WhereOr("rejected != ?", true)
			return q, nil
		}).
		Select(&users)

	if pgutil.IsSqlNoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *storage) getProjectCandidates(ctx context.Context, tags, filteredOutUsers []string) ([]*types.User, error) {
	return nil, nil
}
