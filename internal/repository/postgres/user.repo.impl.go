package postgres

import (
	"context"
	"database/sql"
	"effective-mobile-test-task/internal/apperror"
	"effective-mobile-test-task/internal/model"
	"effective-mobile-test-task/internal/repository"
	"effective-mobile-test-task/internal/types"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/rs/zerolog"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) (repository.UserRepo, error) {
	if db == nil {
		return nil, apperror.NewAppError("NewUserRepo", "db instnce is not initialize", nil)
	}
	return &userRepo{db: db}, nil
}

func (r *userRepo) Find(ctx context.Context, uqo *model.UserQueryOptions) ([]model.User, int, error) {
	log := zerolog.Ctx(ctx).With().Str("method", "userRepo.Find").Logger()
	log.Debug().Interface("userQueryOptions", uqo).Msg("building query for finding users")

	users := make([]model.User, 0)

	if uqo == nil {
		uqo = &model.UserQueryOptions{}
	}
	countBuilder := sq.Select("COUNT(*)").
		From("users").
		PlaceholderFormat(sq.Dollar)

	limit := uqo.GetLimit()
	offset := (uqo.GetPage() - 1) * limit

	builder := sq.Select("uuid", "name", "surname", "patronymic", "age", "gender", "country_id", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		OrderBy(fmt.Sprintf("%s %s", uqo.GetOrderBy(), uqo.GetOrderDir())).
		Offset(offset).
		Limit(limit)

	if uqo.Filter.Surname != nil {
		countBuilder = countBuilder.Where(sq.Eq{"surname": *uqo.Filter.Surname})
		builder = builder.Where(sq.Eq{"surname": *uqo.Filter.Surname})
	}
	if uqo.Filter.Patronymic != nil {
		countBuilder = countBuilder.Where(sq.Eq{"patronymic": *uqo.Filter.Patronymic})
		builder = builder.Where(sq.Eq{"patronymic": *uqo.Filter.Patronymic})
	}
	if uqo.Filter.Age != nil {
		countBuilder = countBuilder.Where(sq.Eq{"age": *uqo.Filter.Age})
		builder = builder.Where(sq.Eq{"age": *uqo.Filter.Age})
	}
	if uqo.Filter.Gender != nil {
		countBuilder = countBuilder.Where(sq.Eq{"gender": *uqo.Filter.Gender})
		builder = builder.Where(sq.Eq{"gender": *uqo.Filter.Gender})
	}
	if uqo.Filter.CountryID != nil {
		countBuilder = countBuilder.Where(sq.Eq{"country_id": *uqo.Filter.CountryID})
		builder = builder.Where(sq.Eq{"country_id": *uqo.Filter.CountryID})
	}

	var totalCount int
	countQuery, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, apperror.NewAppError("userRepo.Find", "failed count sql build", err)
	}

	log.Debug().Str("countQuery", countQuery).Interface("countArgs", countArgs).Msg("executing SQL query for count total rows")
	err = r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, apperror.NewAppError("userRepo.Find", "failed count query", err)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, 0, apperror.NewAppError("userRepo.Find", "failed sql build", err)
	}

	log.Debug().Str("query", query).Interface("args", args).Msg("executing SQL query")
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, apperror.NewAppError("userRepo.Find", "failed query", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.UUID, &u.Name, &u.Surname, &u.Patronymic, &u.Age, &u.Gender, &u.CountryID, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, apperror.NewAppError("userRepo.Find", "failed scan", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, apperror.NewAppError("userRepo.Find", "rows interation error", err)
	}

	log.Debug().Int("page_rows", totalCount).Msg("users found into database")

	return users, totalCount, nil
}

func (r *userRepo) Insert(ctx context.Context, u *model.UserCreate) error {
	log := zerolog.Ctx(ctx).With().Str("method", "userRepo.Insert").Logger()
	log.Debug().Interface("user", u).Msg("starting transaction to insert user")

	query := "INSERT INTO users (uuid, name, surname, patronymic, age, gender, country_id) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	args := []interface{}{u.UUID, u.Name, u.Surname, u.Patronymic, u.Age, u.Gender, u.CountryID}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.NewAppError("userRepo.Insert", "error beginning transaction", err)
	}
	defer tx.Rollback()

	log.Debug().Str("query", query).Interface("args", args).Msg("executing SQL query")
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return apperror.NewAppError("userRepo.Insert", "failed query", err)
	}
	if err = tx.Commit(); err != nil {
		return apperror.NewAppError("userRepo.Insert", "error commiting transaction", err)
	}
	log.Info().Str("uuid", string(u.UUID)).Msg("user inserted into database")

	return nil
}

func (r *userRepo) Update(ctx context.Context, uuid types.UUID, u *model.UserUpdate) (int64, error) {
	log := zerolog.Ctx(ctx).With().Str("method", "userRepo.Update").Logger()
	log.Debug().Interface("user", u).Msg("building query for updating user")

	builder := sq.Update("users").PlaceholderFormat(sq.Dollar).Where(sq.Eq{"uuid": uuid})

	hasUpdates := false

	if u.Name != nil {
		builder = builder.Set("name", *u.Name)
		hasUpdates = true
	}
	if u.Surname != nil {
		builder = builder.Set("surname", *u.Surname)
		hasUpdates = true
	}
	if u.Patronymic != nil {
		builder = builder.Set("patronymic", *u.Patronymic)
		hasUpdates = true
	}
	if u.Age != nil {
		builder = builder.Set("age", *u.Age)
		hasUpdates = true
	}
	if u.Gender != nil {
		builder = builder.Set("gender", *u.Gender)
		hasUpdates = true
	}
	if u.CountryID != nil {
		builder = builder.Set("country_id", *u.CountryID)
		hasUpdates = true
	}

	log.Debug().Bool("hasUpdates", hasUpdates).Msg("checking if any to update")
	if !hasUpdates {
		return 0, apperror.NewAppError("userRepo.Update", "no fields to update", nil)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, apperror.NewAppError("userRepo.Update", "failed build sql", err)
	}

	log.Debug().Str("query", query).Interface("args", args).Msg("executing SQL query")
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, apperror.NewAppError("userRepo.Update", "failed exec", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, apperror.NewAppError("userRepo.Update", "can't get affectedRows count", err)
	}

	log.Info().Int64("affected rows", affected).Msg("user updated into database")

	return affected, nil
}

func (r *userRepo) Delete(ctx context.Context, uuid types.UUID) (int64, error) {
	log := zerolog.Ctx(ctx).With().Str("method", "userRepo.Delete").Logger()

	log.Debug().Str("uuid", string(uuid)).Msg("executing SQL query to delete user")
	result, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE uuid = $1", uuid)
	if err != nil {
		return 0, apperror.NewAppError("userRepo.Delete", "failed exec", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, apperror.NewAppError("userRepo.Delete", "can't get affectedRows count", err)
	}

	log.Info().Int64("affected rows", affected).Msg("user deleted from database")

	return affected, nil
}
