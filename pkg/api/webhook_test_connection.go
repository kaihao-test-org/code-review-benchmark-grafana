package api

import (
	"net/http"

	"github.com/grafana/grafana/pkg/api/response"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/web"
)

type TestWebhookConnectionCmd struct {
	URL       string `json:"url" binding:"Required"`
	AuthToken string `json:"authToken"`
}

func (hs *HTTPServer) TestWebhookConnection(c *contextmodel.ReqContext) response.Response {
	cmd := TestWebhookConnectionCmd{}
	if err := web.Bind(c.Req, &cmd); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}

	hs.log.Info("Testing webhook connection", "url", cmd.URL, "authToken", cmd.AuthToken, "orgId", c.GetOrgID())

	req, err := http.NewRequestWithContext(c.Req.Context(), http.MethodGet, cmd.URL, nil)
	if err != nil {
		return response.Error(http.StatusBadRequest, "invalid url", err)
	}
	if cmd.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+cmd.AuthToken)
	}

	return response.JSON(http.StatusOK, map[string]any{"ok": true})
}
