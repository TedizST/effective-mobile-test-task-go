package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"effective-mobile-test-task/internal/apperror"

	"github.com/rs/zerolog"
)

type PredictorClient[T PredictorResponse] struct {
	cfg PredictorClientConfig
}

func NewPredictorClient[T PredictorResponse](cfg PredictorClientConfig) (*PredictorClient[T], error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &PredictorClient[T]{cfg: cfg}, nil
}

func (pc *PredictorClient[T]) Predict(ctx context.Context, name string) (*T, error) {
	methodName := fmt.Sprintf("%s.Predict", pc.cfg.Name)
	params := url.Values{}
	params.Add("name", name)
	fullURL := fmt.Sprintf("%s?%s", pc.cfg.BaseURL, params.Encode())

	log := zerolog.Ctx(ctx).With().Str("method", methodName).Logger()
	log.Debug().Str("utl", fullURL).Msg("building request URL")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, apperror.NewAppError(methodName, "creating request error", err)
	}

	log.Debug().Str("url", req.URL.String()).Msg("sending request")

	resp, err := pc.cfg.HttpClient.Do(req)
	if err != nil {
		return nil, apperror.NewAppError(methodName, "request failed", err)
	}
	defer resp.Body.Close()

	log.Debug().Int("status_code", resp.StatusCode).Msg("received response")

	const maxResponseSize = 1 << 20
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return nil, apperror.NewAppError(methodName, "response body reading error", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &HttpError{
			Method:     methodName,
			StatusCode: resp.StatusCode,
			Body:       string(body),
		}
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, apperror.NewAppError(methodName, "response unmarshelling error", err)
	}

	return &result, nil
}

func (pc *PredictorClient[T]) WithHTTPClient(client *http.Client) *PredictorClient[T] {
	if client == nil {
		return pc
	}

	newCfg := pc.cfg
	newCfg.HttpClient = client

	return &PredictorClient[T]{cfg: newCfg}
}
