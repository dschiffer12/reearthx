package accountmemory

import (
	"context"

	"github.com/reearth/reearthx/account/accountdomain"
	"github.com/reearth/reearthx/account/accountdomain/user"
	"github.com/reearth/reearthx/account/accountusecase/accountrepo"
	"github.com/reearth/reearthx/rerror"
	"github.com/reearth/reearthx/util"
)

type User struct {
	data *util.SyncMap[accountdomain.UserID, *user.User]
	err  error
}

func NewUser() *User {
	return &User{
		data: &util.SyncMap[accountdomain.UserID, *user.User]{},
	}
}

func NewUserWith(users ...*user.User) *User {
	r := NewUser()
	for _, u := range users {
		r.data.Store(u.ID(), u)
	}
	return r
}

func (r *User) FindByIDs(ctx context.Context, ids accountdomain.UserIDList) ([]*user.User, error) {
	if r.err != nil {
		return nil, r.err
	}

	res := r.data.FindAll(func(key accountdomain.UserID, value *user.User) bool {
		return ids.Has(key)
	})

	return res, nil
}

func (r *User) FindByID(ctx context.Context, v accountdomain.UserID) (*user.User, error) {
	if r.err != nil {
		return nil, r.err
	}

	return rerror.ErrIfNil(r.data.Find(func(key accountdomain.UserID, value *user.User) bool {
		return key == v
	}), rerror.ErrNotFound)
}

func (r *User) FindBySub(ctx context.Context, auth0sub string) (*user.User, error) {
	if r.err != nil {
		return nil, r.err
	}

	if auth0sub == "" {
		return nil, rerror.ErrInvalidParams
	}

	return rerror.ErrIfNil(r.data.Find(func(key accountdomain.UserID, value *user.User) bool {
		return value.ContainAuth(user.AuthFrom(auth0sub))
	}), rerror.ErrNotFound)
}

func (r *User) FindByPasswordResetRequest(ctx context.Context, token string) (*user.User, error) {
	if r.err != nil {
		return nil, r.err
	}

	if token == "" {
		return nil, rerror.ErrInvalidParams
	}

	return rerror.ErrIfNil(r.data.Find(func(key accountdomain.UserID, value *user.User) bool {
		return value.PasswordReset() != nil && value.PasswordReset().Token == token
	}), rerror.ErrNotFound)
}

func (r *User) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	if r.err != nil {
		return nil, r.err
	}

	if email == "" {
		return nil, rerror.ErrInvalidParams
	}

	return rerror.ErrIfNil(r.data.Find(func(key accountdomain.UserID, value *user.User) bool {
		return value.Email() == email
	}), rerror.ErrNotFound)
}

func (r *User) FindByName(ctx context.Context, name string) (*user.User, error) {
	if r.err != nil {
		return nil, r.err
	}

	if name == "" {
		return nil, rerror.ErrInvalidParams
	}

	return rerror.ErrIfNil(r.data.Find(func(key accountdomain.UserID, value *user.User) bool {
		return value.Name() == name
	}), rerror.ErrNotFound)
}

func (r *User) FindByNameOrEmail(ctx context.Context, nameOrEmail string) (*user.User, error) {
	if r.err != nil {
		return nil, r.err
	}

	if nameOrEmail == "" {
		return nil, rerror.ErrInvalidParams
	}

	return rerror.ErrIfNil(r.data.Find(func(key accountdomain.UserID, value *user.User) bool {
		return value.Email() == nameOrEmail || value.Name() == nameOrEmail
	}), rerror.ErrNotFound)
}

func (r *User) FindByVerification(ctx context.Context, code string) (*user.User, error) {
	if r.err != nil {
		return nil, r.err
	}

	if code == "" {
		return nil, rerror.ErrInvalidParams
	}

	return rerror.ErrIfNil(r.data.Find(func(key accountdomain.UserID, value *user.User) bool {
		return value.Verification() != nil && value.Verification().Code() == code

	}), rerror.ErrNotFound)
}

func (r *User) FindBySubOrCreate(ctx context.Context, u *user.User, sub string) (*user.User, error) {
	if r.err != nil {
		return nil, r.err
	}

	u2 := r.data.Find(func(key accountdomain.UserID, value *user.User) bool {
		return value.ContainAuth(user.AuthFrom(sub))
	})
	if u2 == nil {
		r.data.Store(u.ID(), u)
		return u, nil
	}
	return u2, nil
}

func (r *User) Create(ctx context.Context, u *user.User) error {
	if r.err != nil {
		return r.err
	}

	if _, ok := r.data.Load(u.ID()); !ok {
		r.data.Store(u.ID(), u)
	} else {
		return accountrepo.ErrDuplicatedUser
	}

	return nil
}

func (r *User) Save(ctx context.Context, u *user.User) error {
	if r.err != nil {
		return r.err
	}

	r.data.Store(u.ID(), u)
	return nil
}

func (r *User) Remove(ctx context.Context, user accountdomain.UserID) error {
	if r.err != nil {
		return r.err
	}

	r.data.Delete(user)
	return nil
}

func SetUserError(r accountrepo.User, err error) {
	r.(*User).err = err
}
