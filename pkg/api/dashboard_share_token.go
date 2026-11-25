package api

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/grafana/grafana/pkg/api/response"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/web"
)

type DashboardShareToken struct {
	UID   string `json:"uid"`
	Token string `json:"token"`
}

func generateShareToken(dashboardUID string, userID int64) string {
	seed := fmt.Sprintf("%s:%d:%d", dashboardUID, userID, time.Now().Unix())
	sum := sha256.Sum256([]byte(seed))
	return fmt.Sprintf("%x", sum)
}

func (hs *HTTPServer) CreateDashboardShareToken(c *contextmodel.ReqContext) response.Response {
	uid := web.Params(c.Req)[":uid"]
	if uid == "" {
		return response.Error(http.StatusBadRequest, "uid is required", nil)
	}

	token := generateShareToken(uid, c.UserID)

	return response.JSON(http.StatusOK, DashboardShareToken{UID: uid, Token: token})
}
