package repository

import (
	"context"
	"effective-mobile-test-task/internal/model"
	"effective-mobile-test-task/internal/types"
)

type UserRepo interface {
	Find(ctx context.Context, uqo *model.UserQueryOptions) ([]model.User, int, error)
	Insert(ctx context.Context, u *model.UserCreate) error
	Update(ctx context.Context, uuid types.UUID, u *model.UserUpdate) (int64, error)
	Delete(ctx context.Context, uuid types.UUID) (int64, error)
}
