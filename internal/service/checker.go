package service

import (
	"context"
	"net/http"
	"time"

	"github.com/Dokhoyan/2025-11-12-test/internal/domain"
)

type LinkChecker interface {
	CheckLink(ctx context.Context, url string) (domain.LinkStatus, error)
	CheckLinks(ctx context.Context, urls []string) ([]domain.Link, error)
}

type HTTPLinkChecker struct {
	client  *http.Client
	timeout time.Duration
}

func NewHTTPLinkChecker(timeout time.Duration) *HTTPLinkChecker {
	return &HTTPLinkChecker{
		client: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

func (c *HTTPLinkChecker) CheckLink(ctx context.Context, url string) (domain.LinkStatus, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return domain.StatusUnavailable, err
	}

	req.Header.Set("User-Agent", "LinkChecker/1.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return domain.StatusUnavailable, nil
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return domain.StatusAvailable, nil
	}

	return domain.StatusUnavailable, nil
}

func (c *HTTPLinkChecker) CheckLinks(ctx context.Context, urls []string) ([]domain.Link, error) {
	type result struct {
		link  domain.Link
		index int
	}

	results := make([]domain.Link, len(urls))
	resultChan := make(chan result, len(urls))

	for i, url := range urls {
		go func(idx int, u string) {
			status, _ := c.CheckLink(ctx, u)
			resultChan <- result{
				link: domain.Link{
					URL:    u,
					Status: string(status),
				},
				index: idx,
			}
		}(i, url)
	}

	for i := 0; i < len(urls); i++ {
		res := <-resultChan
		results[res.index] = res.link
	}

	return results, nil
}
