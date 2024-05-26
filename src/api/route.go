package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	RouteRegisterPartner = "register_partner"
	RouteUploadAvatar    = "upload_avatar"
)

type RouteInfo = struct {
	Path        string
	Method      string
	Handler     func(c *gin.Context)
	RequireAuth bool
}

func (s *Server) AllRoutes() map[string]RouteInfo {
	return map[string]RouteInfo{
		RouteRegisterPartner: {
			Path:        "/partner/register",
			Method:      http.MethodPost,
			Handler:     s.RegisterPartner,
			RequireAuth: false,
		},
		RouteUploadAvatar: {
			Path:        "/user/avatar/upload",
			Method:      http.MethodPost,
			Handler:     s.HandleUploadAvatar,
			RequireAuth: false,
		},
	}
}
