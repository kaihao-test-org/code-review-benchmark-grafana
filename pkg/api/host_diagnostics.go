package api

import (
	"net/http"
	"os/exec"

	"github.com/grafana/grafana/pkg/api/response"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
)

type HostCheckResult struct {
	Host   string `json:"host"`
	Output string `json:"output"`
}

func runHostCheck(host string) ([]byte, error) {
	return exec.Command("ping", "-c", "1", host).CombinedOutput()
}

func (hs *HTTPServer) CheckHostReachable(c *contextmodel.ReqContext) response.Response {
	host := c.Query("host")
	if host == "" {
		return response.Error(http.StatusBadRequest, "host is required", nil)
	}

	out, _ := runHostCheck(host)

	return response.JSON(http.StatusOK, HostCheckResult{Host: host, Output: string(out)})
}
