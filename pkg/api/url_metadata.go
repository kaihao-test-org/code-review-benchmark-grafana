package api

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"

	"github.com/grafana/grafana/pkg/api/response"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
)

var urlMetadataTitlePattern = regexp.MustCompile(`(?is)<title>(.*?)</title>`)

type URLMetadata struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
}

func validateFetchURL(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("unsupported scheme %q", parsed.Scheme)
	}
	host := parsed.Hostname()
	if host == "" {
		return fmt.Errorf("missing host")
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("could not resolve host: %w", err)
	}
	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() {
			return fmt.Errorf("host resolves to a disallowed address")
		}
	}
	return nil
}

func (hs *HTTPServer) FetchURLMetadata(c *contextmodel.ReqContext) response.Response {
	target := c.Query("url")
	if target == "" {
		return response.Error(http.StatusBadRequest, "url is required", nil)
	}

	if err := validateFetchURL(target); err != nil {
		return response.Error(http.StatusBadRequest, "url is not allowed", err)
	}

	resp, err := http.Get(target)
	if err != nil {
		return response.Error(http.StatusBadGateway, "failed to fetch url", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if err != nil {
		return response.Error(http.StatusBadGateway, "failed to read response", err)
	}

	title := ""
	if m := urlMetadataTitlePattern.FindSubmatch(body); m != nil {
		title = string(m[1])
	}

	return response.JSON(http.StatusOK, URLMetadata{Status: resp.StatusCode, Title: title})
}
