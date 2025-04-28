package handler

import (
	"effective-mobile-test-task/internal/apperror"
	"effective-mobile-test-task/internal/dto"
	"effective-mobile-test-task/internal/model"
	"effective-mobile-test-task/internal/service"
	"effective-mobile-test-task/internal/types"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) (*UserHandler, error) {
	if userService == nil {
		return nil, apperror.NewAppError("NewUserHandler", "userService is required", nil)
	}

	return &UserHandler{userService: userService}, nil
}

func (uh *UserHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", uh.FindUsers)
	r.Post("/", uh.CreateUser)
	r.Patch("/{uuid}", uh.UpdateUser)
	r.Delete("/{uuid}", uh.DeleteUser)

	return r
}

// FindUsers godoc
// @Summary Поиск пользователей с использованием фильтров и пагинации
// @Description Получение списка пользователей с помощью передачи query параметров
// @Tags users
// @Accept json
// @Produce json
// @Param page query int true "Номер страницы"
// @Param limit query int true "Количество записей на 1 странице"
// @Param name query string false "Имя пользователя"
// @Param surname query string false "Фамилия пользователя"
// @Param patronymic query string false "Отчество пользователя"
// @Param age query int false "Возраст пользователя"
// @Param gender query string false "Пол пользователя"
// @Param country_id query string false "Код страны пользователя"
// @Param order_by query string false "Поле для сортировки"
// @Param order_dir query string false "Направление сортировки ASC, DESC"
// @Success 200 {object} dto.ResponseDTO{payload=dto.ListOfUsersPayload}
// @Failure 400 {object} dto.ErrorResponseDTO
// @Failure 500 {object} dto.ErrorResponseDTO
// @Router /users [get]
func (uh *UserHandler) FindUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uqo := &model.UserQueryOptions{}

	if page := r.FormValue("page"); page != "" {
		uint64Page, err := strconv.ParseUint(page, 10, 64)
		if err != nil {
			errorResponse(ctx, w, apperror.NewHttpError(400, "page must be a positive number"))
			return
		}
		uqo.Page = types.Page(uint64Page)
	} else {
		errorResponse(ctx, w, apperror.NewHttpError(400, "page is required"))
		return
	}
	if limit := r.FormValue("limit"); limit != "" {
		uint64Limit, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			errorResponse(ctx, w, apperror.NewHttpError(400, "limit must be a positive number"))
			return
		}
		uqo.Limit = types.Limit(uint64Limit)
	} else {
		errorResponse(ctx, w, apperror.NewHttpError(400, "limit is required"))
		return
	}

	if name := r.FormValue("name"); name != "" {
		uqo.Filter.Name = (*types.Name)(&name)
	}
	if surname := r.FormValue("surname"); surname != "" {
		uqo.Filter.Surname = (*types.Surname)(&surname)
	}
	if patronymic := r.FormValue("patronymic"); patronymic != "" {
		uqo.Filter.Patronymic = (*types.Patronymic)(&patronymic)
	}
	if age := r.FormValue("age"); age != "" {
		uintAge, err := strconv.ParseUint(age, 10, 0)
		if err != nil {
			errorResponse(ctx, w, apperror.NewHttpError(400, "age must be a positive number"))
			return
		}
		uqo.Filter.Age = (*types.Age)(&uintAge)
	}
	if gender := r.FormValue("gender"); gender != "" {
		uqo.Filter.Gender = (*types.Gender)(&gender)
	}
	if countryID := r.FormValue("country_id"); countryID != "" {
		uqo.Filter.CountryID = (*types.CountryID)(&countryID)
	}

	if orderBy := r.FormValue("order_by"); orderBy != "" {
		if uqo.IsValidOrderBy(orderBy) {
			uqo.OrderBy = (*types.OrderBy)(&orderBy)
		} else {
			errorResponse(ctx, w, apperror.NewHttpError(400, "orderBy has invalid value"))
			return
		}
	}
	if orderDir := r.FormValue("order_dir"); orderDir != "" {
		if uqo.IsValidOrderDir(orderDir) {
			uqo.OrderDir = (*types.OrderDir)(&orderDir)
		} else {
			errorResponse(ctx, w, apperror.NewHttpError(400, "orderDir has invalid value"))
			return
		}
	}

	usersList, err := uh.userService.FindUsers(ctx, uqo)
	if err != nil {
		errorResponse(ctx, w, err)
		return
	}

	successResponse(ctx, w, 200, usersList)
}

// CreateUser godoc
// @Summary Создание пользователя
// @Description Создание пользователя, данные будут обогащены с помощью публичных API
// @Tags users
// @Accept json
// @Produce json
// @Param name body string true "Имя пользователя"
// @Param surname body string true "Фамилия пользователя"
// @Param patronymic query string false "Отчество пользователя"
// @Success 200 {object} dto.ResponseDTO{payload=dto.UserCreatePayload}
// @Failure 400 {object} dto.ErrorResponseDTO
// @Failure 500 {object} dto.ErrorResponseDTO
// @Router /users [post]
func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ucDTO := &dto.UserCreateDTO{}
	err := json.NewDecoder(r.Body).Decode(ucDTO)
	if err != nil {
		errorResponse(ctx, w, apperror.NewHttpError(400, "invalid user json structure"))
		return
	}

	if ucDTO.Name == "" {
		errorResponse(ctx, w, apperror.NewHttpError(400, "name is empty"))
		return
	}
	if ucDTO.Surname == "" {
		errorResponse(ctx, w, apperror.NewHttpError(400, "surname is empty"))
		return
	}

	uuid, err := uh.userService.CreateUser(ctx, ucDTO)
	if err != nil {
		errorResponse(ctx, w, err)
		return
	}

	successResponse(ctx, w, 200, &dto.UserCreatePayload{UUID: uuid})
}

// UpdateUser godoc
// @Summary Обновление данных пользователя
// @Description Обновление данных пользователя (в теле запроса нет обязательных полей, но в случае передачи пустого тела запроса будет возвращен ответ с кодом 400)
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "ID пользователя"
// @Param name body string false "Имя пользователя"
// @Param surname body string false "Фамилия пользователя"
// @Param patronymic body string false "Отчество пользователя"
// @Param age body int false "Возраст пользователя"
// @Param gender body string false "Пол пользователя"
// @Param country_id body string false "Код страны пользователя"
// @Success 200 {object} dto.EmptyResponseDTO
// @Failure 400 {object} dto.ErrorResponseDTO
// @Failure 404 {object} dto.ErrorResponseDTO
// @Failure 500 {object} dto.ErrorResponseDTO
// @Router /users/{uuid} [patch]
func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	uuid := r.PathValue("uuid")
	if uuid == "" {
		errorResponse(ctx, w, apperror.NewHttpError(400, "uuid is empty"))
		return
	}

	uuDTO := &dto.UserUpdateDTO{}
	err := json.NewDecoder(r.Body).Decode(uuDTO)
	if err != nil {
		errorResponse(ctx, w, apperror.NewHttpError(400, "invalid user json structure"))
		return
	}

	hasUpdates := uuDTO.Name != nil || uuDTO.Surname != nil || uuDTO.Patronymic != nil ||
		uuDTO.Age != nil || uuDTO.Gender != nil || uuDTO.CountryID != nil
	if !hasUpdates {
		errorResponse(ctx, w, apperror.NewHttpError(400, "no payload provided for update"))
		return
	}

	err = uh.userService.UpdateUser(ctx, types.UUID(uuid), uuDTO)
	if err != nil {
		errorResponse(ctx, w, err)
		return
	}

	successResponse(ctx, w, 200, nil)
}

// DeleteUser godoc
// @Summary Удаление пользователя
// @Description Удаление пользователя по переданному ID
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "ID пользователя"
// @Success 200 {object} dto.EmptyResponseDTO
// @Failure 400 {object} dto.ErrorResponseDTO
// @Failure 404 {object} dto.ErrorResponseDTO
// @Failure 500 {object} dto.ErrorResponseDTO
// @Router /users/{uuid} [delete]
func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	uuid := r.PathValue("uuid")
	if uuid == "" {
		errorResponse(ctx, w, apperror.NewHttpError(400, "uuid is empty"))
		return
	}

	err := uh.userService.DeleteUser(ctx, types.UUID(uuid))
	if err != nil {
		errorResponse(ctx, w, err)
		return
	}

	successResponse(ctx, w, 200, nil)
}
