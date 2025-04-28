package configs

import (
	"effective-mobile-test-task/internal/httpclient"
	"fmt"
	"os"
	"strings"
)

func GetAgifyConfig() (*httpclient.PredictorClientConfig, error) {
	token := os.Getenv("API_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("API_TOKEN is required")
	}

	baseURL := os.Getenv("AGIFY_BASE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("AGIFY_BASE_URL is required")
	}

	return &httpclient.PredictorClientConfig{
		Name:    httpclient.Agify,
		Token:   token,
		BaseURL: strings.TrimSuffix(baseURL, "/"),
	}, nil
}
