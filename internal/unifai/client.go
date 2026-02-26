package unifai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	endpoint string
	apiKey   string
	http     *http.Client
}

func NewClient(endpoint, apiKey string, timeout time.Duration) *Client {
	return &Client{
		endpoint: strings.TrimRight(endpoint, "/"),
		apiKey:   apiKey,
		http: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) Search(ctx context.Context, req SearchRequest) (any, error) {
	query := url.Values{}
	query.Set("query", req.Query)
	query.Set("limit", fmt.Sprintf("%d", req.Limit))
	query.Set("offset", fmt.Sprintf("%d", req.Offset))
	if len(req.IncludeActions) > 0 {
		query.Set("includeActions", strings.Join(req.IncludeActions, ","))
	}

	return c.doJSON(ctx, http.MethodGet, "/actions/search", query, nil)
}

func (c *Client) Invoke(ctx context.Context, req InvokeRequest) (any, error) {
	body := map[string]any{
		"action": req.Action,
	}
	if req.Payload != nil {
		body["payload"] = req.Payload
	}
	if req.Payment != nil {
		body["payment"] = req.Payment
	}

	return c.doJSON(ctx, http.MethodPost, "/actions/call", nil, body)
}

func (c *Client) doJSON(ctx context.Context, method, path string, query url.Values, requestBody any) (any, error) {
	fullURL := c.endpoint + path
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	var bodyReader io.Reader
	if requestBody != nil {
		data, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", c.apiKey)
	httpReq.Header.Set("Accept", "application/json")
	if requestBody != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: strings.TrimSpace(string(data))}
	}

	if len(data) == 0 {
		return map[string]any{}, nil
	}

	var decoded any
	if err := json.Unmarshal(data, &decoded); err != nil {
		return nil, fmt.Errorf("decode response body: %w", err)
	}

	return decoded, nil
}
