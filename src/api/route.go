package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	RouteRegisterPartner = "register_partner"
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
	}
}
