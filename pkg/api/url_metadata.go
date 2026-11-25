package api

import (
	"io"
	"net/http"
	"regexp"

	"github.com/grafana/grafana/pkg/api/response"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
)

var urlMetadataTitlePattern = regexp.MustCompile(`(?is)<title>(.*?)</title>`)

type URLMetadata struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
}

func (hs *HTTPServer) FetchURLMetadata(c *contextmodel.ReqContext) response.Response {
	target := c.Query("url")
	if target == "" {
		return response.Error(http.StatusBadRequest, "url is required", nil)
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
