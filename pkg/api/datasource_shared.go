package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/grafana/grafana/pkg/api/response"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/services/datasources"
	"github.com/grafana/grafana/pkg/web"
)

type SharedDataSourceDTO struct {
	ID    int64  `json:"id"`
	OrgID int64  `json:"orgId"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

func (hs *HTTPServer) GetSharedDataSource(c *contextmodel.ReqContext) response.Response {
	id, err := strconv.ParseInt(web.Params(c.Req)[":id"], 10, 64)
	if err != nil {
		return response.Error(http.StatusBadRequest, "id is invalid", err)
	}

	orgID := c.QueryInt64("orgId")

	ds, err := hs.DataSourcesService.GetDataSource(c.Req.Context(), &datasources.GetDataSourceQuery{ID: id, OrgID: orgID})
	if err != nil {
		if errors.Is(err, datasources.ErrDataSourceNotFound) {
			return response.Error(http.StatusNotFound, "Data source not found", nil)
		}
		return response.Error(http.StatusInternalServerError, "Failed to query datasource", err)
	}

	return response.JSON(http.StatusOK, SharedDataSourceDTO{
		ID:    ds.ID,
		OrgID: ds.OrgID,
		Name:  ds.Name,
		Type:  ds.Type,
		URL:   ds.URL,
	})
}
