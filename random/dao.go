package random

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"nobl9-recruitment-task/logging"
	"nobl9-recruitment-task/rest"
)

func NewDao() *dao {
	var config struct {
		Timeout time.Duration `default:"10s"`
		Url     string        `required:"true"`
	}
	envconfig.MustProcess("RANDOM_INTEGERS_GENERATOR", &config)
	logging.Info("creating new random generator service dao with URL: %s and timeout: %s", config.Url, config.Timeout.String())
	return &dao{
		client:  http.Client{},
		url:     config.Url,
		timeout: config.Timeout,
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
	timeout          time.Duration
	predefinedParams map[string]string
}

func (d dao) GetIntegers(ctx context.Context, length int) ([]int, error) {
	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, d.url, nil)
	if err != nil {
		return nil, err
	}
	query := url.Values{}
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
		return nil, rest.NewHttpError(req.URL.Path, string(raw), resp.StatusCode)
	}
	split := strings.Split(strings.TrimSpace(string(raw)), "\n")
	res := make([]int, len(split))
	for i := range split {
		converted, err := strconv.Atoi(split[i])
		if err != nil {
			return nil, rest.NewHttpError(req.URL.Path, "malformed response from random generator, one of the integers was not an integer al all", http.StatusInternalServerError)
		}
		res[i] = converted
	}
	return res, nil
}
