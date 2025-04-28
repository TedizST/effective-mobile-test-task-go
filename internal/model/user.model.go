package model

import (
	"effective-mobile-test-task/internal/types"
	"time"
)

type (
	UserCreate struct {
		UUID       types.UUID
		Name       types.Name
		Surname    types.Surname
		Patronymic *types.Patronymic
		Age        *types.Age
		Gender     *types.Gender
		CountryID  *types.CountryID
	}
	UserUpdate struct {
		Name       *types.Name
		Surname    *types.Surname
		Patronymic *types.Patronymic
		Age        *types.Age
		Gender     *types.Gender
		CountryID  *types.CountryID
	}
	User struct {
		UUID       types.UUID
		Name       types.Name
		Surname    types.Surname
		Patronymic *types.Patronymic
		Age        *types.Age
		Gender     *types.Gender
		CountryID  *types.CountryID
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}
	UserFilter struct {
		Name       *types.Name
		Surname    *types.Surname
		Patronymic *types.Patronymic
		Age        *types.Age
		Gender     *types.Gender
		CountryID  *types.CountryID
	}
	UserQueryOptions struct {
		Filter   UserFilter
		OrderBy  *types.OrderBy
		OrderDir *types.OrderDir
		Pagination
	}
)

const (
	UUID       = "uuid"
	Name       = "name"
	Surname    = "surname"
	Patronymic = "patronymic"
	Age        = "age"
	Gender     = "gender"
	CountryId  = "country_id"
	CreatedAt  = "created_at"
	ASC        = "ASC"
	DESC       = "DESC"
)

func (uqo *UserQueryOptions) IsValidOrderBy(field string) bool {
	switch field {
	case UUID, Name, Surname, Patronymic, Age, Gender, CountryId, CreatedAt:
		return true
	}
	return false
}

func (uqo *UserQueryOptions) IsValidOrderDir(direction string) bool {
	switch direction {
	case ASC, DESC:
		return true
	}
	return false
}

func (uqo *UserQueryOptions) GetOrderBy() string {
	if uqo.OrderBy == nil {
		return string(CreatedAt)
	}
	return string(*uqo.OrderBy)
}

func (uqo *UserQueryOptions) GetOrderDir() string {
	if uqo.OrderDir == nil {
		return DESC
	}
	return string(*uqo.OrderDir)
}
