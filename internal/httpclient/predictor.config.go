package httpclient

import (
	"fmt"
	"net/http"
	"time"
)

type PredictorClientConfig struct {
	Name       APIType
	Token      string
	BaseURL    string
	Timeout    time.Duration
	HttpClient *http.Client
}

func (c *PredictorClientConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("API name is required")
	}
	if c.Token == "" {
		return fmt.Errorf("API token is required")
	}
	if c.BaseURL == "" {
		return fmt.Errorf("base URL is required")
	}
	if c.Timeout == 0 {
		c.Timeout = 10 * time.Second
	}
	if c.HttpClient == nil {
		c.HttpClient = &http.Client{Timeout: c.Timeout}
	}
	return nil
}
