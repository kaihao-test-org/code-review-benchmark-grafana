package api

import (
	"net/http"

	"github.com/grafana/grafana/pkg/api/response"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/util"
	"github.com/grafana/grafana/pkg/web"
)

type DashboardShareToken struct {
	UID   string `json:"uid"`
	Token string `json:"token"`
}

func generateShareToken(dashboardUID string, userID int64) (string, error) {
	return util.GetRandomString(32)
}

func (hs *HTTPServer) CreateDashboardShareToken(c *contextmodel.ReqContext) response.Response {
	uid := web.Params(c.Req)[":uid"]
	if uid == "" {
		return response.Error(http.StatusBadRequest, "uid is required", nil)
	}

	token, err := generateShareToken(uid, c.UserID)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "Failed to generate token", err)
	}

	return response.JSON(http.StatusOK, DashboardShareToken{UID: uid, Token: token})
}
