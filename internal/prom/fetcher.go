package prom

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Client struct {
	BaseURL string
	hc      *http.Client
}

func NewClient() *Client {
	base := os.Getenv("PROM_URL")
	if base == "" {
		base = "http://localhost:9090"
	}
	return &Client{
		BaseURL: base,
		hc:      &http.Client{Timeout: 20 * time.Second},
	}
}

func (c *Client) ListMetrics(ctx context.Context, limit int) ([]string, error) {
	u := c.BaseURL + "/api/v1/label/__name__/values"
	req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
	res, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("prometheus error: %d", res.StatusCode)
	}
	var out struct {
		Status string   `json:"status"`
		Data   []string `json:"data"`
	}
	b, _ := io.ReadAll(res.Body)
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	if limit > 0 && len(out.Data) > limit {
		return out.Data[:limit], nil
	}
	return out.Data, nil
}

func (c *Client) InstantQuery(ctx context.Context, promql string) ([]byte, error) {
	q := url.QueryEscape(promql)
	u := fmt.Sprintf("%s/api/v1/query?query=%s", c.BaseURL, q)
	req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
	res, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("prometheus error: %d", res.StatusCode)
	}
	return io.ReadAll(res.Body)
}

