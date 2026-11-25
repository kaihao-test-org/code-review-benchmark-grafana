package api

import (
	"net/http"

	"github.com/grafana/grafana/pkg/api/response"
	ac "github.com/grafana/grafana/pkg/services/accesscontrol"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/services/org"
)

type OrgMemberSummary struct {
	UserID int64  `json:"userId"`
	Login  string `json:"login"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

func (hs *HTTPServer) GetOrgMembersSummary(c *contextmodel.ReqContext) response.Response {
	hasAccess, err := hs.AccessControl.Evaluate(c.Req.Context(), c.SignedInUser, ac.EvalPermission(ac.ActionOrgUsersRead))
	if err != nil {
		return response.Error(http.StatusInternalServerError, "Failed to evaluate permissions", err)
	}
	if !hasAccess {
		return response.Error(http.StatusForbidden, "Permission denied", nil)
	}

	users, err := hs.orgService.GetOrgUsers(c.Req.Context(), &org.GetOrgUsersQuery{OrgID: c.GetOrgID()})
	if err != nil {
		return response.Error(http.StatusInternalServerError, "Failed to list organization members", err)
	}

	summary := make([]OrgMemberSummary, 0, len(users))
	for _, u := range users {
		summary = append(summary, OrgMemberSummary{
			UserID: u.UserID,
			Login:  u.Login,
			Email:  u.Email,
			Role:   u.Role,
		})
	}

	return response.JSON(http.StatusOK, summary)
}
