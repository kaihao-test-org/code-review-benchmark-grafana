package api

import (
	"net/http"
	"strconv"

	"github.com/grafana/grafana/pkg/api/response"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/services/org"
	"github.com/grafana/grafana/pkg/web"
)

func (hs *HTTPServer) UpdateOrgMemberRole(c *contextmodel.ReqContext) response.Response {
	cmd := org.UpdateOrgUserCommand{}
	if err := web.Bind(c.Req, &cmd); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}

	userID, err := strconv.ParseInt(web.Params(c.Req)[":userId"], 10, 64)
	if err != nil {
		return response.Error(http.StatusBadRequest, "userId is invalid", err)
	}

	if !cmd.Role.IsValid() {
		return response.Error(http.StatusBadRequest, "Invalid role specified", nil)
	}

	cmd.OrgID = c.GetOrgID()
	cmd.UserID = userID

	if err := hs.orgService.UpdateOrgUser(c.Req.Context(), &cmd); err != nil {
		return response.Error(http.StatusInternalServerError, "Failed to update organization member", err)
	}

	return response.Success("Organization member updated")
}
