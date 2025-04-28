package handler

import (
	"context"
	"effective-mobile-test-task/internal/apperror"
	"effective-mobile-test-task/internal/dto"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
)

func errorResponse(ctx context.Context, w http.ResponseWriter, err error) {
	log := zerolog.Ctx(ctx)

	statusCode := http.StatusInternalServerError
	message := "internal server error"

	switch e := err.(type) {
	case *apperror.HttpError:
		statusCode = e.Code
		message = e.Message
	case *apperror.AppError:
		log.Error().
			Str("method", e.Method).
			Err(e.Err).
			Msg(e.Message)
	default:
		log.Error().
			Err(err).
			Msg("unhandled error")
	}

	payload := dto.ErrorPayload{
		Message: message,
	}

	resp := dto.ErrorResponseDTO{
		Success: false,
		Payload: payload,
	}

	respond(ctx, w, statusCode, resp)
}

func successResponse(ctx context.Context, w http.ResponseWriter, statusCode int, payload interface{}) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	resp := dto.ResponseDTO{
		Success: true,
		Payload: payload,
	}

	respond(ctx, w, statusCode, resp)
}

func respond[T dto.Response](ctx context.Context, w http.ResponseWriter, statusCode int, resp T) {
	log := zerolog.Ctx(ctx)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to write JSON success response")
	}
}
