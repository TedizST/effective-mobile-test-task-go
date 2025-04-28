package configs

import (
	"effective-mobile-test-task/internal/httpclient"
	"fmt"
	"os"
	"strings"
)

func GetGenderizeConfig() (*httpclient.PredictorClientConfig, error) {
	token := os.Getenv("API_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("API_TOKEN is required")
	}

	baseURL := os.Getenv("GENDERIZE_BASE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("GENDERIZE_BASE_URL is required")
	}

	return &httpclient.PredictorClientConfig{
		Name:    httpclient.Genderize,
		Token:   token,
		BaseURL: strings.TrimSuffix(baseURL, "/"),
	}, nil
}
