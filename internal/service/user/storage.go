package user

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/osapers/mch-back/internal/constant"
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

func (s *storage) contextWithDBSecret(ctx context.Context) context.Context {
	return context.WithValue(ctx, constant.DBSecretCtxKey, s.conn.Secret)
}

func (s *storage) getByID(ctx context.Context, id string) (*types.User, error) {
	u := &types.User{ID: id}

	err := s.conn.DB.ModelContext(s.contextWithDBSecret(ctx), u).WherePK().First()

	if pgutil.IsSqlNoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *storage) getByEmail(ctx context.Context, email string) (*types.User, error) {
	u := &types.User{}

	err := s.conn.DB.ModelContext(s.contextWithDBSecret(ctx), u).Where("email = ?", email).First()

	if pgutil.IsSqlNoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *storage) create(ctx context.Context, email, password string) (*types.User, error) {
	u := &types.User{Email: email, Password: password}

	_, err := s.conn.DB.ModelContext(s.contextWithDBSecret(ctx), u).Returning("*").Insert()
	err = pgutil.HandleSqlIntegrityViolation(err)

	return u, err
}

func (s *storage) getParticipatedEvents(ctx context.Context, userID string) (map[string]struct{}, error) {
	var events []string

	err := s.conn.DB.ModelContext(ctx, (*types.UserToEvent)(nil)).
		Column("event_id").
		Where("user_id = ?", userID).
		Select(&events)

	if pgutil.IsSqlNoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	res := make(map[string]struct{})
	for _, e := range events {
		res[e] = struct{}{}
	}

	return res, nil
}

func (s *storage) participateInEvent(ctx context.Context, userID, eventID string) error {
	_, err := s.conn.DB.ModelContext(ctx, &types.UserToEvent{UserID: userID, EventID: eventID}).Insert()
	err = pgutil.HandleSqlIntegrityViolation(err)

	return err
}

func (s *storage) update(ctx context.Context, user *types.User) (*types.User, error) {
	_, err := s.conn.DB.ModelContext(s.contextWithDBSecret(ctx), user).WherePK().Returning("*").UpdateNotZero()
	err = pgutil.HandleSqlIntegrityViolation(err)
	return user, err
}

func (s *storage) getAppliedProjects(ctx context.Context, userID string) (map[string]*types.UserToProject, error) {
	var projects []*types.UserToProject

	err := s.conn.DB.ModelContext(ctx, &projects).
		Where("user_id = ?", userID).
		Where("applied = ?", true).
		Select(&projects)

	if pgutil.IsSqlNoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	res := make(map[string]*types.UserToProject)

	for _, p := range projects {
		res[p.ProjectID] = p
	}

	return res, nil
}

func (s *storage) getViewedProjects(ctx context.Context, userID string) ([]string, error) {
	var projects []string

	err := s.conn.DB.ModelContext(ctx, (*types.UserToProject)(nil)).
		Column("project_id").
		Where("user_id = ?", userID).
		Select(&projects)

	if pgutil.IsSqlNoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *storage) viewProject(ctx context.Context, userID, projectID string) error {
	_, err := s.conn.DB.ModelContext(ctx, &types.UserToProject{UserID: userID, ProjectID: projectID}).Insert()
	err = pgutil.HandleSqlIntegrityViolation(err)

	return err
}

func (s *storage) applyToProject(ctx context.Context, userID, projectID string) error {
	_, err := s.conn.DB.ModelContext(ctx, &types.UserToProject{UserID: userID, ProjectID: projectID, Applied: true}).WherePK().UpdateNotZero()
	err = pgutil.HandleSqlIntegrityViolation(err)

	return err
}

func (s *storage) deleteProject(ctx context.Context, userID, projectID string) error {
	_, err := s.conn.DB.ModelContext(ctx, &types.Project{ID: projectID}).WherePK().Where("owner_id = ?", userID).Delete()

	if pgutil.IsSqlNoRows(err) || err == nil {
		return nil
	}

	return err
}

func (s *storage) getMyProjects(ctx context.Context, userID string, appliedIds []string) ([]*types.Project, error) {
	var projects []*types.Project

	q := s.conn.DB.ModelContext(ctx, &projects).WhereOr("owner_id = ?", userID)

	if len(appliedIds) != 0 {
		q = q.WhereOr("id in (?)", pg.In(appliedIds))
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

func (s *storage) getOwnersInfo(ctx context.Context, projectOwnerIds []string) (map[string]*types.User, error) {
	var owners []*types.User

	q := s.conn.DB.ModelContext(s.contextWithDBSecret(ctx), &owners)

	if len(projectOwnerIds) > 0 {
		q = q.Where("id in (?)", pg.In(projectOwnerIds))
	}

	err := q.Select()

	if pgutil.IsSqlNoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	res := make(map[string]*types.User)

	for _, owner := range owners {
		res[owner.ID] = owner
	}

	return res, nil
}

func (s *storage) matchTags(ctx context.Context, keywords []string) ([]string, error) {
	var tags []*types.Tag

	err := s.conn.DB.ModelContext(ctx, &tags).
		WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			for _, kw := range keywords {
				q = q.WhereOr("name ILIKE ?", kw)
			}
			return q, nil
		}).Select(&tags)

	if pgutil.IsSqlNoRows(err) {
		return []string{}, nil
	}

	if err != nil {
		return nil, err
	}

	var res []string

	for _, tag := range tags {
		res = append(res, tag.Name)
	}

	return res, nil
}
