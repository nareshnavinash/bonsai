package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// Client is an OpenAI-compatible HTTP client that talks to llama-server
// or any endpoint implementing the OpenAI chat completions API.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a client pointing at the given base URL (without /v1).
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

// NewClientFromEnv creates a client using BONSAI_HOST env var,
// defaulting to http://127.0.0.1:8081.
func NewClientFromEnv() *Client {
	host := os.Getenv("BONSAI_HOST")
	if host == "" {
		host = DefaultBaseURL()
	}
	return NewClient(host)
}

// DefaultBaseURL returns the default server URL based on BONSAI_PORT or 8081.
func DefaultBaseURL() string {
	port := os.Getenv("BONSAI_PORT")
	if port == "" {
		port = "8081"
	}
	return "http://127.0.0.1:" + port
}

// Health checks if the server is healthy by calling GET /health.
func (c *Client) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/health", nil)
	if err != nil {
		return err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("server not reachable: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server unhealthy: status %d", resp.StatusCode)
	}
	return nil
}

// Models calls GET /v1/models and returns the available models.
func (c *Client) Models(ctx context.Context) (*ModelList, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/v1/models", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot list models: %w", err)
	}
	defer resp.Body.Close()

	var result ModelList
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("invalid models response: %w", err)
	}
	return &result, nil
}

// ChatCompletion sends a non-streaming chat completion request.
func (c *Client) ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	req.Stream = false
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("chat request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Error.Message != "" {
			return nil, fmt.Errorf("chat error: %s", errResp.Error.Message)
		}
		return nil, fmt.Errorf("chat request failed: status %d", resp.StatusCode)
	}

	var result ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("invalid chat response: %w", err)
	}
	return &result, nil
}

// ChatCompletionStream sends a streaming chat completion request.
// The callback is called for each SSE chunk containing content.
func (c *Client) ChatCompletionStream(ctx context.Context, req *ChatCompletionRequest, fn func(ChatCompletionResponse) error) error {
	req.Stream = true
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("chat request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Error.Message != "" {
			return fmt.Errorf("chat error: %s", errResp.Error.Message)
		}
		return fmt.Errorf("chat request failed: status %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")

		if data == "[DONE]" {
			break
		}

		var chunk ChatCompletionResponse
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue // skip malformed chunks
		}

		if err := fn(chunk); err != nil {
			return err
		}
	}

	return scanner.Err()
}
