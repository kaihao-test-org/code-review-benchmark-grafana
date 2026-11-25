package api

import (
	"errors"
	"net/http"

	"github.com/grafana/grafana/pkg/api/response"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/services/datasources"
	"github.com/grafana/grafana/pkg/web"
)

type DataSourceExportDTO struct {
	UID            string            `json:"uid"`
	Name           string            `json:"name"`
	Type           string            `json:"type"`
	URL            string            `json:"url"`
	User           string            `json:"user"`
	SecureJsonData map[string]string `json:"secureJsonData"`
}

func buildDataSourceExportDTO(ds *datasources.DataSource, secrets map[string]string) DataSourceExportDTO {
	return DataSourceExportDTO{
		UID:            ds.UID,
		Name:           ds.Name,
		Type:           ds.Type,
		URL:            ds.URL,
		User:           ds.User,
		SecureJsonData: secrets,
	}
}

func (hs *HTTPServer) ExportDataSource(c *contextmodel.ReqContext) response.Response {
	ds, err := hs.getRawDataSourceByUID(c.Req.Context(), web.Params(c.Req)[":uid"], c.GetOrgID())
	if err != nil {
		if errors.Is(err, datasources.ErrDataSourceNotFound) {
			return response.Error(http.StatusNotFound, "Data source not found", nil)
		}
		return response.Error(http.StatusInternalServerError, "Failed to query datasource", err)
	}

	secrets, err := hs.DataSourcesService.DecryptedValues(c.Req.Context(), ds)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "Failed to decrypt datasource secrets", err)
	}

	return response.JSON(http.StatusOK, buildDataSourceExportDTO(ds, secrets))
}
