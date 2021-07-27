package random

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func NewDao() *dao {
	return &dao{
		client: http.Client{},
		//url:    "https://www.random.org/integers",
		url: "http://localhost:8081/integers",
		predefinedParams: map[string]string{
			"min":    "1",
			"max":    "100000",
			"col":    "1",
			"base":   "10",
			"format": "plain",
			"rnd":    "new",
		},
	}
}

type dao struct {
	client           http.Client
	url              string
	predefinedParams map[string]string
}

func (d dao) GetIntegers(ctx context.Context, length int) ([]int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, d.url, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("num", strconv.Itoa(length))
	for k, v := range d.predefinedParams {
		query.Set(k, v)
	}
	req.URL.RawQuery = query.Encode()
	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response from: %s with code: %d and message: %s", d.url, resp.StatusCode, string(raw))
	}
	split := strings.Split(strings.TrimSpace(string(raw)), "\n")
	res := make([]int, len(split))
	for i := range split {
		converted, err := strconv.Atoi(split[i])
		if err != nil {
			return nil, fmt.Errorf("malformed response from: %s, one of the integers was not an integer at all", d.url)
		}
		res[i] = converted
	}
	return res, nil
}
