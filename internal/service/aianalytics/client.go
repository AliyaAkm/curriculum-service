package aianalytics

import (
	"context"
	"curriculum-service/internal/domain/studentstats"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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
		return nil, nil
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 5 * time.Second
	}

	return &Client{
		baseURL: baseURL,
		http:    &http.Client{Timeout: cfg.Timeout},
	}, nil
}

type analyticsResponse struct {
	UserID            uuid.UUID         `json:"userId"`
	TotalRequests     int               `json:"totalRequests"`
	CompletedRequests int               `json:"completedRequests"`
	FailedRequests    int               `json:"failedRequests"`
	Chats             int               `json:"chats"`
	UserMessages      int               `json:"userMessages"`
	AssistantMessages int               `json:"assistantMessages"`
	InputTokens       int               `json:"inputTokens"`
	OutputTokens      int               `json:"outputTokens"`
	AvgLatencyMS      int               `json:"avgLatencyMs"`
	LastActivityAt    *time.Time        `json:"lastActivityAt"`
	DailyActivity     []dailyAIActivity `json:"dailyActivity"`
}

type dailyAIActivity struct {
	Date     string `json:"date"`
	Requests int    `json:"requests"`
}

func (c *Client) GetUserAnalytics(ctx context.Context, userID uuid.UUID) (*studentstats.AIStats, []studentstats.ActivityDay, error) {
	if c == nil || c.baseURL == "" {
		return nil, nil, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/analytics/users/"+userID.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, nil, fmt.Errorf("ai analytics status: %s", resp.Status)
	}

	var data analyticsResponse
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, nil, err
	}

	activity := make([]studentstats.ActivityDay, len(data.DailyActivity))
	for i := range data.DailyActivity {
		activity[i] = studentstats.ActivityDay{
			Date:       data.DailyActivity[i].Date,
			AIRequests: data.DailyActivity[i].Requests,
		}
	}

	return &studentstats.AIStats{
		Available:         true,
		TotalRequests:     data.TotalRequests,
		CompletedRequests: data.CompletedRequests,
		FailedRequests:    data.FailedRequests,
		Chats:             data.Chats,
		UserMessages:      data.UserMessages,
		AssistantMessages: data.AssistantMessages,
		InputTokens:       data.InputTokens,
		OutputTokens:      data.OutputTokens,
		AvgLatencyMS:      data.AvgLatencyMS,
		LastActivityAt:    data.LastActivityAt,
	}, activity, nil
}
