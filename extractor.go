// Package traefik_url_extractor a Traefik plugin to extract URL path parameters into request headers.
package traefik_url_extractor

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
)

// Config the plugin configuration.
type Config struct {
	Regex   string            `json:"regex"`
	Headers map[string]string `json:"headers,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Headers: make(map[string]string),
	}
}

// URLExtractor extracts URL path parameters into request headers.
type URLExtractor struct {
	next    http.Handler
	re      *regexp.Regexp
	headers map[string]string
	name    string
}

// New creates a new URLExtractor plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.Regex == "" {
		return nil, fmt.Errorf("regex cannot be empty")
	}
	if len(config.Headers) == 0 {
		return nil, fmt.Errorf("headers cannot be empty")
	}

	re, err := regexp.Compile(config.Regex)
	if err != nil {
		return nil, fmt.Errorf("invalid regex: %w", err)
	}

	return &URLExtractor{
		next:    next,
		re:      re,
		headers: config.Headers,
		name:    name,
	}, nil
}

func (u *URLExtractor) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	matches := u.re.FindStringSubmatch(req.URL.String())
	if matches == nil {
		u.next.ServeHTTP(rw, req)
		return
	}

	groups := make(map[string]string)
	for i, name := range u.re.SubexpNames() {
		if name != "" && i < len(matches) {
			groups[name] = matches[i]
		}
	}

	for headerName, groupName := range u.headers {
		if val, ok := groups[groupName]; ok {
			req.Header.Set(headerName, val)
		}
	}

	u.next.ServeHTTP(rw, req)
}
