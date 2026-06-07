package coderunner

import (
	"bytes"
	"context"
	"curriculum-service/internal/usecase/codeattempt"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ClientConfig struct {
	BaseURL string
	Timeout time.Duration
}

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(cfg ClientConfig) (*Client, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if baseURL == "" {
		return nil, fmt.Errorf("code runner service url is required")
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 10 * time.Second
	}

	return &Client{
		baseURL: baseURL,
		http:    &http.Client{Timeout: cfg.Timeout},
	}, nil
}

type runRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type runResponse struct {
	Output string `json:"output"`
	Error  string `json:"error"`
	Passed bool   `json:"passed"`
}

func (c *Client) Run(ctx context.Context, language, code string) (codeattempt.RunnerResult, error) {
	payload, err := json.Marshal(runRequest{
		Language: language,
		Code:     code,
	})
	if err != nil {
		return codeattempt.RunnerResult{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/run", bytes.NewReader(payload))
	if err != nil {
		return codeattempt.RunnerResult{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return codeattempt.RunnerResult{}, err
	}
	defer resp.Body.Close()

	var result runResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return codeattempt.RunnerResult{}, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		if result.Error != "" {
			return codeattempt.RunnerResult(result), fmt.Errorf("code runner error: %s", result.Error)
		}
		return codeattempt.RunnerResult(result), fmt.Errorf("code runner status: %s", resp.Status)
	}

	return codeattempt.RunnerResult(result), nil
}
