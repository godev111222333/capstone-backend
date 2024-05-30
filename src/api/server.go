package api

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/godev111222333/capstone-backend/src/misc"
	"github.com/godev111222333/capstone-backend/src/store"
	"github.com/godev111222333/capstone-backend/src/token"
)

const DefaultHost = "0.0.0.0"

type Server struct {
	cfg          *misc.ApiServerConfig
	route        *gin.Engine
	store        *store.DbStore
	s3store      *store.S3Store
	tokenMaker   token.Maker
	hashVerifier *misc.HashVerifier
	otpService   *OTPService
}

func NewServer(
	cfg *misc.ApiServerConfig,
	store *store.DbStore,
	s3Store *store.S3Store,
	otpService *OTPService,
) *Server {
	route := gin.New()
	tokenMaker, err := token.NewJWTMaker("12345678901234567890123456789012")
	if err != nil {
		panic(err)
	}

	server := &Server{
		cfg, route, store, s3Store, tokenMaker, misc.NewHashVerifier(), otpService,
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
	s.route.Use(cors.Default())
}

func (s *Server) registerHandlers() {
	authGroup := s.route.Group("/").Use(authMiddleware(s.tokenMaker))

	for _, r := range s.AllRoutes() {
		if !r.RequireAuth {
			s.route.Handle(r.Method, r.Path, r.Handler)
		} else {
			authGroup.Handle(r.Method, r.Path, r.Handler)
		}
	}
}
