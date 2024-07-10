package api

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/godev111222333/capstone-backend/src/misc"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/store"
	"github.com/godev111222333/capstone-backend/src/token"
)

const DefaultHost = "0.0.0.0"

const (
	ChanBufferSize = 100
)

type Server struct {
	cfg            *misc.ApiServerConfig
	feCfg          *misc.FEConfig
	route          *gin.Engine
	store          *store.DbStore
	s3store        *store.S3Store
	tokenMaker     token.Maker
	hashVerifier   *misc.HashVerifier
	otpService     *OTPService
	pdfService     IPDFService
	paymentService IPaymentService

	bankMetadata []string
	chatRooms    sync.Map

	adminSubs sync.Map

	adminNotificationQueue    chan NotificationMsg
	adminNewConversationQueue chan ConversationMsg
}

func NewServer(
	cfg *misc.ApiServerConfig,
	feCfg *misc.FEConfig,
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
		feCfg,
		route,
		store,
		s3Store,
		tokenMaker,
		misc.NewHashVerifier(),
		otpService,
		pdfService,
		paymentService,
		bankMetadata,
		sync.Map{},
		sync.Map{},
		make(chan NotificationMsg, ChanBufferSize),
		make(chan ConversationMsg, ChanBufferSize),
	}
	server.setUp()
	return server
}

func (s *Server) Run() error {
	fmt.Printf("API server running at port: %s\n", s.cfg.ApiPort)
	s.startAdminSub()

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
					adminGroup.Handle(r.Method, strings.TrimPrefix(r.Path, "/admin/"), r.Handler)
				case model.RoleNamePartner:
					partnerGroup.Handle(r.Method, strings.TrimPrefix(r.Path, "/partner/"), r.Handler)
				case model.RoleNameCustomer:
					customerGroup.Handle(r.Method, strings.TrimPrefix(r.Path, "/customer/"), r.Handler)
				}
			}
		}
	}
}
