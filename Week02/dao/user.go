package dao

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

type User struct {
	Id int64
}

func NewUser(id int64) *User {
	return &User{Id: id}
}

func query() (*User, error) {
	return &User{}, sql.ErrNoRows
}

func (u *User) Get(ctx context.Context) (*User, error) {
	user, err := query()
	if err == sql.ErrNoRows {
		return nil, errors.WithMessagef(err, "uid: %d", u.Id)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "get user failed, params: [uid: %d]", u.Id)
	}
	return user, nil
}
