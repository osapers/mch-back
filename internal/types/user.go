package types

import (
	"context"

	"github.com/osapers/mch-back/internal/constant"
	"github.com/osapers/mch-back/pkg/crypt"
)

type User struct {
	tableName        struct{} `pg:"mch.user,alias:u"` // nolint:structcheck,unused
	ID               string   `json:"id" pg:"type:uuid"`
	FirstName        string   `json:"first_name"`
	LastName         string   `json:"last_name"`
	MiddleName       string   `json:"middle_name"`
	Photo            string   `json:"photo"`
	About            string   `json:"about"`
	Email            string   `json:"email"`
	Phone            string   `json:"phone"`
	Password         string   `json:"-"`
	Tags             []string `json:"tags" pg:",array"`
	TagsIntersection []string `json:"tags_intersection" pg:"-"`
}

func (u *User) ToJson() *User {
	user := *u
	if len(user.Tags) == 0 {
		user.Tags = []string{}
	}

	return &user
}

// BeforeInsert executes before sql insert
func (u *User) BeforeInsert(ctx context.Context) (context.Context, error) {
	var err error

	if u.Password != "" {
		if u.Password, err = crypt.Encrypt(u.Password, ctx.Value(constant.DBSecretCtxKey).([]byte)); err != nil {
			return nil, err
		}
	}

	return ctx, nil
}

// BeforeUpdate executes before sql update
func (u *User) BeforeUpdate(ctx context.Context) (context.Context, error) {
	return u.BeforeInsert(ctx)
}

// AfterScan executes after sql select
func (u *User) AfterScan(ctx context.Context) error {
	var err error

	if u.Password != "" {
		if u.Password, err = crypt.Decrypt(u.Password, ctx.Value(constant.DBSecretCtxKey).([]byte)); err != nil {
			return err
		}
	}

	return nil
}
