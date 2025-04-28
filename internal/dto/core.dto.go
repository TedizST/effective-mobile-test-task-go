package dto

type (
	// ResponseDTO общее представление ответа
	ResponseDTO struct {
		Success bool        `json:"success" default:"true"` // Статус операции
		Payload interface{} `json:"payload,omitempty"`      // Полезная нагрузка
	}
	// EmptyResponseDTO в ответе присутствует только статус
	EmptyResponseDTO struct {
		Success bool `json:"success" default:"true"` // Статус операции
	}
	// ErrorPayload полезная нагрузка с текстом ошибки
	ErrorPayload struct {
		Message string `json:"message"` // Текст ошибки
	}
	// ErrorResponseDTO ответ при ошибке
	ErrorResponseDTO struct {
		Success bool `json:"success" default:"false"`
		Payload ErrorPayload
	}
	Response interface {
		ResponseDTO | ErrorResponseDTO
	}
)
