package dto

import (
	"effective-mobile-test-task/internal/types"
)

type (
	// UserCreateDTO представляет данные для создания пользователя в БД
	UserCreateDTO struct {
		Name       types.Name        `json:"name" example:"Dmitriy"`                    // Имя пользователя
		Surname    types.Surname     `json:"surname" example:"Ushakov"`                 // Фамилия пользователя
		Patronymic *types.Patronymic `json:"patronymic,omitempty" example:"Vasilevich"` // Отчество пользователя (необязательное поле)
	}
	// UserUpdateDTO поля для обнолвения данных пользователя
	UserUpdateDTO struct {
		Name       *types.Name       `json:"name,omitempty" example:"Dmitriy"`                         // Имя пользователя
		Surname    *types.Surname    `json:"surname,omitempty" example:"Ushakov"`                      // Фамилия пользователя
		Patronymic *types.Patronymic `json:"patronymic,omitempty" example:"Vasilevich"`                // Отчество пользователя
		Age        *types.Age        `json:"age,omitempty" example:"22"`                               // Возврат пользователя
		Gender     *types.Gender     `json:"gender,omitempty" example:"male"`                          // Пол пользователя
		CountryID  *types.CountryID  `json:"country_id,omitempty" example:"2006-01-02T15:04:05Z07:00"` // Строковый ID страны пользователя
	}
	// UserPayload поля для обнолвения данных пользователя
	UserPayload struct {
		UUID       types.UUID        `json:"uuid" example:"8d571787-9981-4add-a713-2fde6236e84b"` // ID пользователя
		Name       types.Name        `json:"name" example:"Dmitriy"`                              // Имя пользователя
		Surname    types.Surname     `json:"surname" example:"Ushakov"`                           // Фамилия пользователя
		Patronymic *types.Patronymic `json:"patronymic,omitempty" example:"Vasilevich"`           // Отчество пользователя
		Age        *types.Age        `json:"age,omitempty" example:"22"`                          // Возврат пользователя
		Gender     *types.Gender     `json:"gender,omitempty" example:"male"`                     // Пол пользователя
		CountryID  *types.CountryID  `json:"country_id,omitempty" example:"RU"`                   // Строковый ID страны пользователя
		CreatedAt  string            `json:"created_at" example:"2006-01-02T15:04:05Z07:00"`      // Строковое представление даты создания пользователя
	}
	// ListOfUsersPayload полезная нагрузка со списком пользователей
	ListOfUsersPayload struct {
		Total int           `json:"total" example:"0"` // Общее количество записей с переданными фильтрами
		Users []UserPayload `json:"users"`             // Список пользователей на указанной странице
	}
	// UserCreatePayload полезная нагрузка, содержащая информацию о созданном ползователе
	UserCreatePayload struct {
		UUID types.UUID `json:"uuid" example:"8d571787-9981-4add-a713-2fde6236e84b"` // ID пользователя
	}
)
