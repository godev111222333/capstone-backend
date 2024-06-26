package api

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/godev111222333/capstone-backend/src/model"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/godev111222333/capstone-backend/src/misc"
	"github.com/godev111222333/capstone-backend/src/store"
	"github.com/godev111222333/capstone-backend/src/token"
)

const DefaultHost = "0.0.0.0"

type Server struct {
	cfg            *misc.ApiServerConfig
	route          *gin.Engine
	store          *store.DbStore
	s3store        *store.S3Store
	tokenMaker     token.Maker
	hashVerifier   *misc.HashVerifier
	otpService     *OTPService
	pdfService     IPDFService
	paymentService IPaymentService

	bankMetadata []string
}

func NewServer(
	cfg *misc.ApiServerConfig,
	store *store.DbStore,
	s3Store *store.S3Store,
	otpService *OTPService,
	bankMetadata []string,
	pdfService IPDFService,
	paymentService IPaymentService,
) *Server {
	route := gin.New()
	tokenMaker, err := token.NewJWTMaker("12345678901234567890123456789012")
	if err != nil {
		panic(err)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for tagName, valFunc := range CustomValidations {
			if err := v.RegisterValidation(tagName, valFunc); err != nil {
				panic(err)
			}
		}
	}

	server := &Server{
		cfg,
		route,
		store,
		s3Store,
		tokenMaker,
		misc.NewHashVerifier(),
		otpService,
		pdfService,
		paymentService,
		bankMetadata,
	}
	server.setUp()
	return server
}

func (s *Server) Run() error {
	fmt.Printf("API server running at port: %s\n", s.cfg.ApiPort)

	return s.route.Run(fmt.Sprintf("%s:%s", DefaultHost, s.cfg.ApiPort))
}

func (s *Server) setUp() {
	s.registerMiddleware()
	s.registerHandlers()
}

func (s *Server) registerMiddleware() {
	s.route.Use(customCORSHeader())
}

func (s *Server) registerHandlers() {
	authGroup := s.route.Group("/").Use(authMiddleware(s.tokenMaker), s.activeAccountMiddleware())
	adminGroup := s.route.Group("/admin").Use(authMiddleware(s.tokenMaker), s.activeAccountMiddleware(), s.authRole(model.RoleNameAdmin))
	partnerGroup := s.route.Group("/partner").Use(authMiddleware(s.tokenMaker), s.activeAccountMiddleware(), s.authRole(model.RoleNamePartner))
	customerGroup := s.route.Group("/customer").Use(authMiddleware(s.tokenMaker), s.activeAccountMiddleware(), s.authRole(model.RoleNameCustomer))
	for _, r := range s.AllRoutes() {
		if !r.RequireAuth {
			s.route.Handle(r.Method, r.Path, r.Handler)
		} else {
			if len(r.AuthRoles) == 0 {
				authGroup.Handle(r.Method, r.Path, r.Handler)
				continue
			}

			for _, authRole := range r.AuthRoles {
				switch authRole {
				case model.RoleNameAdmin:
					truncatedPrefix := strings.Replace(r.Path, "/admin", "", 1)
					adminGroup.Handle(r.Method, truncatedPrefix, r.Handler)
				case model.RoleNamePartner:
					truncatedPrefix := strings.Replace(r.Path, "/partner", "", 1)
					partnerGroup.Handle(r.Method, truncatedPrefix, r.Handler)
				case model.RoleNameCustomer:
					truncatedPrefix := strings.Replace(r.Path, "/customer", "", 1)
					customerGroup.Handle(r.Method, truncatedPrefix, r.Handler)
				}
			}
		}
	}
}

func customCORSHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
