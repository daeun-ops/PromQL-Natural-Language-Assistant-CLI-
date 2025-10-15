package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Client struct {
	apiKey string
	model  string
}

func New() *Client {
	return &Client{
		apiKey: os.Getenv("OPENAI_API_KEY"),
		model:  envOr("OPENAI_MODEL", "gpt-4.1-mini"),
	}
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// GeneratePromQL asks the LLM to return a JSON object: {"promql": "..."}
func (c *Client) GeneratePromQL(ctx context.Context, userRequest string, schemaHint string) (string, error) {
	// Fallback to mock if no key
	if c.apiKey == "" {
		return mockPromQL(userRequest), nil
	}

	system := `You are a PromQL generator. 
Given a natural-language request and a list of available metrics/labels, output ONLY a strict JSON: {"promql": "<expression>"}.
- Do not add explanations.
- Prefer histogram_quantile on *_bucket when user asks for percentiles.
- Use only metrics/labels that likely exist in the given schema hint. 
- Keep time windows from the user's text (default to 5m if missing).
- Avoid expensive queries (no 30d lookbacks).`

	user := fmt.Sprintf(`Available schema hint (truncated and optional):
%s

User request:
%s

Return only: {"promql": "<expression>"}`, schemaHint, userRequest)

	// Minimal Responses API payload (compatible with OpenAI Responses)
	payload := map[string]any{
		"model": c.model,
		"input": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": user},
		},
		"response_format": map[string]string{"type": "json_object"},
		"temperature":     0.2,
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/responses", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Timeout: 45 * time.Second}
	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return "", fmt.Errorf("openai error: status %d", res.StatusCode)
	}

	var out map[string]any
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return "", err
	}

	// Responses API returns an "output" array with objects that contain "content"
	// Each content item may have "text" (JSON string). We'll try to parse it.
	promql, err := extractPromQL(out)
	if err != nil {
		return "", err
	}
	return promql, nil
}

func extractPromQL(res map[string]any) (string, error) {
	// Try to find a text block containing a JSON with {"promql": "..."}
	// Shape: output[0].content[0].text
	output, ok := res["output"].([]any)
	if !ok || len(output) == 0 {
		return "", errors.New("invalid response: missing output")
	}
	first := output[0].(map[string]any)
	content, ok := first["content"].([]any)
	if !ok || len(content) == 0 {
		return "", errors.New("invalid response: missing content")
	}
	// find first text field
	for _, c := range content {
		m, _ := c.(map[string]any)
		if txt, ok := m["text"].(string); ok && strings.TrimSpace(txt) != "" {
			// txt should be a JSON like {"promql":"..."}
			var v struct{ PromQL string `json:"promql"` }
			if err := json.Unmarshal([]byte(txt), &v); err == nil && v.PromQL != "" {
				return v.PromQL, nil
			}
			// Sometimes model may wrap code blocks; try to strip
			trim := strings.TrimSpace(txt)
			trim = strings.Trim(trim, "`")
			if err := json.Unmarshal([]byte(trim), &v); err == nil && v.PromQL != "" {
				return v.PromQL, nil
			}
		}
	}
	return "", errors.New("failed to parse promql from response")
}

func mockPromQL(input string) string {
	in := strings.ToLower(input)
	switch {
	case strings.Contains(in, "에러") || strings.Contains(in, "error") || strings.Contains(in, "5xx"):
		return `sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))`
	case strings.Contains(in, "p95") || strings.Contains(in, "latency") || strings.Contains(in, "지연"):
		return `histogram_quantile(0.95, sum by (le, handler)(rate(prometheus_http_request_duration_seconds_bucket[5m])))`
	case strings.Contains(in, "업") && strings.Contains(in, "타겟") || strings.Contains(in, "targets"):
		return `sum(up)`
	default:
		return `sum(rate(http_requests_total[5m]))`
	}
}

