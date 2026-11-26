package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/grafana/grafana/pkg/api/response"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/util"
)

func (hs *HTTPServer) templateExportDir() string {
	return filepath.Join(hs.Cfg.DataPath, "exports")
}

func (hs *HTTPServer) DownloadTemplateExport(c *contextmodel.ReqContext) response.Response {
	name := c.Query("name")
	if name == "" {
		return response.Error(http.StatusBadRequest, "name is required", nil)
	}

	cleaned, err := util.CleanRelativePath(name)
	if err != nil {
		return response.Error(http.StatusBadRequest, "invalid name", err)
	}

	path := filepath.Join(hs.templateExportDir(), cleaned)
	data, err := os.ReadFile(path)
	if err != nil {
		return response.Error(http.StatusNotFound, "export not found", err)
	}

	return response.CreateNormalResponse(
		http.Header{"Content-Type": {"application/octet-stream"}},
		data,
		http.StatusOK,
	)
}
