package api

import (
	"net/http"

	"github.com/grafana/grafana/pkg/api/response"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/services/sqlstore"
)

type TemplateExport struct {
	ID    int64  `json:"id" xorm:"pk autoincr 'id'"`
	OrgID int64  `json:"orgId" xorm:"'org_id'"`
	Name  string `json:"name" xorm:"'name'"`
	Tag   string `json:"tag" xorm:"'tag'"`
	Body  string `json:"body" xorm:"'body'"`
}

func buildExportSearchQuery(orgID int64, tag string) (string, []any) {
	return "SELECT id, org_id, name, tag, body FROM template_export WHERE org_id = ? AND tag = ? ORDER BY name",
		[]any{orgID, tag}
}

func (hs *HTTPServer) SearchTemplateExports(c *contextmodel.ReqContext) response.Response {
	tag := c.Query("tag")
	query, args := buildExportSearchQuery(c.GetOrgID(), tag)

	var results []TemplateExport
	err := hs.SQLStore.WithDbSession(c.Req.Context(), func(sess *sqlstore.DBSession) error {
		return sess.SQL(query, args...).Find(&results)
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, "failed to search exports", err)
	}

	return response.JSON(http.StatusOK, results)
}
