package service

import (
	"context"
	"effective-mobile-test-task/internal/apperror"
	"effective-mobile-test-task/internal/dto"
	"effective-mobile-test-task/internal/httpclient"
	"effective-mobile-test-task/internal/model"
	"effective-mobile-test-task/internal/repository"
	"effective-mobile-test-task/internal/types"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type UserService struct {
	userRepo          repository.UserRepo
	agifyClient       *httpclient.PredictorClient[httpclient.AgifyResponse]
	genderizeClient   *httpclient.PredictorClient[httpclient.GenderizeResponse]
	nationalizeClient *httpclient.PredictorClient[httpclient.NationalizeResponse]
}

func NewUserService(
	userRepo repository.UserRepo,
	agifyClient *httpclient.PredictorClient[httpclient.AgifyResponse],
	genderizeClient *httpclient.PredictorClient[httpclient.GenderizeResponse],
	nationalizeClient *httpclient.PredictorClient[httpclient.NationalizeResponse]) (*UserService, error) {
	methodName := "NewUserService"

	if userRepo == nil {
		return nil, apperror.NewAppError(methodName, "userRepo is required", nil)
	}
	if agifyClient == nil {
		return nil, apperror.NewAppError(methodName, "agifyClient is required", nil)
	}
	if genderizeClient == nil {
		return nil, apperror.NewAppError(methodName, "genderizeClient is required", nil)
	}
	if nationalizeClient == nil {
		return nil, apperror.NewAppError(methodName, "nationalizeClient is required", nil)
	}

	return &UserService{
		userRepo:          userRepo,
		agifyClient:       agifyClient,
		genderizeClient:   genderizeClient,
		nationalizeClient: nationalizeClient,
	}, nil
}

func (us *UserService) FindUsers(ctx context.Context, uqo *model.UserQueryOptions) (*dto.ListOfUsersPayload, error) {
	log := zerolog.Ctx(ctx).With().Str("method", "UserService.FindUsers").Logger()

	log.Debug().Interface("userQueryOptions", uqo).Msg("extracting users with filters from database")
	users, total, err := us.userRepo.Find(ctx, uqo)
	if err != nil {
		return nil, err
	}

	log.Debug().Msg("converting []models.User list to []dto.UserResponseDTO")
	usersDTO := []dto.UserPayload{}
	for _, u := range users {
		usersDTO = append(usersDTO, dto.UserPayload{
			UUID:       u.UUID,
			Name:       u.Name,
			Surname:    u.Surname,
			Patronymic: u.Patronymic,
			Age:        u.Age,
			Gender:     u.Gender,
			CountryID:  u.CountryID,
			CreatedAt:  u.CreatedAt.Format(time.RFC3339),
		})
	}

	log.Info().Int("total", total).Int("on_page", len(usersDTO)).Msg("users found in service")

	return &dto.ListOfUsersPayload{
		Total: total,
		Users: usersDTO,
	}, nil
}

func (us *UserService) CreateUser(ctx context.Context, uDTO *dto.UserCreateDTO) (types.UUID, error) {
	log := zerolog.Ctx(ctx).With().Str("method", "UserService.CreateUser").Logger()
	log.Debug().Interface("userDTO", uDTO).Msg("received update user request")

	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	uuidStr := uuid.String()
	log.Debug().Str("uuid", uuidStr).Msg("generated uuid")

	u := &model.UserCreate{
		UUID:       types.UUID(uuidStr),
		Name:       uDTO.Name,
		Surname:    uDTO.Surname,
		Patronymic: uDTO.Patronymic,
	}
	log.Debug().Str("uuid", uuidStr).Interface("user", u).Msg("converted DTO into model")

	nameStr := string(u.Name)
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		log.Debug().Str("uuid", uuidStr).Str("name", nameStr).Msg("calling agify API")
		res, err := us.agifyClient.Predict(ctx, nameStr)
		if err != nil {
			log.Warn().Str("uuid", uuidStr).Str("name", nameStr).Err(err).Msg("failed to call agify")
			return
		}

		mu.Lock()
		u.Age = (*types.Age)(&res.Age)
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		log.Debug().Str("uuid", uuidStr).Str("name", nameStr).Msg("calling genderize API")
		res, err := us.genderizeClient.Predict(ctx, nameStr)
		if err != nil {
			log.Warn().Str("uuid", uuidStr).Str("name", nameStr).Err(err).Msg("failed to call genderize")
			return
		}

		mu.Lock()
		u.Gender = (*types.Gender)(&res.Gender)
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		log.Debug().Str("uuid", uuidStr).Str("name", nameStr).Msg("calling nationalize API")
		res, err := us.nationalizeClient.Predict(ctx, nameStr)
		if err != nil {
			log.Warn().Str("uuid", uuidStr).Str("name", nameStr).Err(err).Msg("failed to call nationalize")
			return
		}
		if len(res.Countries) == 0 {
			log.Warn().Str("uuid", uuidStr).Str("name", nameStr).Msg("propability country list is empty")
			return
		}

		mu.Lock()
		u.CountryID = (*types.CountryID)(&res.Countries[0].CountryId)
		mu.Unlock()
	}()

	wg.Wait()

	log.Debug().Str("uuid", uuidStr).Interface("user", u).Msg("enhanced model and inserting user into database")
	err = us.userRepo.Insert(ctx, u)
	if err != nil {
		return "", err
	}

	log.Info().Str("uuid", uuidStr).Msg("user successfully created")

	return u.UUID, nil
}

func (us *UserService) UpdateUser(ctx context.Context, uuid types.UUID, uDTO *dto.UserUpdateDTO) error {
	log := zerolog.Ctx(ctx).With().Str("method", "UserService.UpdateUser").Logger()

	log.Debug().Str("uuid", string(uuid)).Interface("userDTO", uDTO).Msg("received update user request")
	u := &model.UserUpdate{
		Name:       uDTO.Name,
		Surname:    uDTO.Surname,
		Patronymic: uDTO.Patronymic,
		Age:        uDTO.Age,
		Gender:     uDTO.Gender,
		CountryID:  uDTO.CountryID,
	}
	log.Debug().Str("uuid", string(uuid)).Interface("user", u).Msg("converted dto to model")

	affected, err := us.userRepo.Update(ctx, uuid, u)
	if err != nil {
		return err
	}

	log.Debug().Str("uuid", string(uuid)).Int64("affected", affected).Msg("received response from db")
	if affected == 0 {
		return apperror.NewHttpError(404, "user not found")
	}

	log.Info().Str("uuid", string(uuid)).Msg("user updated succesfully")
	return nil
}

func (us *UserService) DeleteUser(ctx context.Context, uuid types.UUID) error {
	log := zerolog.Ctx(ctx).With().Str("method", "UserService.DeleteUser").Logger()

	log.Debug().Str("uuid", string(uuid)).Msg("received delete user request")

	affected, err := us.userRepo.Delete(ctx, uuid)
	if err != nil {
		return err
	}

	log.Debug().Str("uuid", string(uuid)).Int64("affected", affected).Msg("received response from db")
	if affected == 0 {
		return apperror.NewHttpError(404, "user not found")
	}

	log.Info().Str("uuid", string(uuid)).Msg("user deleted succesfully")
	return nil
}
