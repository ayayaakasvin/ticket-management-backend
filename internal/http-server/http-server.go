package httpserver

import (
	"net/http"
	"sync"
	"time"

	"github.com/ayayaakasvin/oneflick-ticket/internal/config"
	"github.com/ayayaakasvin/oneflick-ticket/internal/http-server/handlers"
	"github.com/ayayaakasvin/oneflick-ticket/internal/http-server/middlewares"
	"github.com/ayayaakasvin/oneflick-ticket/internal/models/inner"

	"github.com/ayayaakasvin/lightmux"
	"github.com/sirupsen/logrus"
)

type ServerApp struct {
	server 		*http.Server

	lmux 		*lightmux.LightMux

	authRepo 	inner.UserRepository
	cache 		inner.Cache

	cfg      	*config.HTTPServer
	wg       	*sync.WaitGroup

	logger 		*logrus.Logger
}

func NewServerApp(cfg *config.HTTPServer, logger *logrus.Logger, wg *sync.WaitGroup, authRepo inner.UserRepository, cache inner.Cache) *ServerApp {
	return &ServerApp{
		cfg:      	cfg,
		logger:   	logger,
		wg:       	wg,
		authRepo: 	authRepo,
		cache: 		cache,
	}
}

func (s *ServerApp) Run() {
	defer s.wg.Done()

	s.setupServer()

	s.setupLightMux()

	s.startServer()
}

func (s *ServerApp) startServer() {
	s.logger.Infof("Server has been started on port: %s", s.cfg.Address)
	s.logger.Infof("Available handlers:\n")

	s.lmux.PrintMiddlewareInfo()
	s.lmux.PrintRoutes()

	go func() {
		ticker := time.NewTicker(time.Minute * 5)
		for range ticker.C {
			s.logger.Info("Server is running...")
		}
	}()

	// RunTLS can be run when server is hosted on domain, acts as seperate service of file storing, for my project, id chose to encapsulate servers under one docker-compose and make nginx-gateaway for my api like auth, file, user service
	// if err := s.lmux.RunTLS(s.cfg.TLS.CertFile, s.cfg.TLS.KeyFile); err != nil {
	if err := s.lmux.Run(); err != nil {
		s.logger.Fatalf("Server exited with error: %v", err)
	}
}

// setuping server by pointer, so we dont have to return any value
func (s *ServerApp) setupServer() {
	if s.server == nil {
		// s.logger.Warn("Server is nil, creating a new server pointer")
		s.server = &http.Server{}
	}

	s.server.Addr = s.cfg.Address
	s.server.IdleTimeout = s.cfg.IdleTimeout
	s.server.ReadTimeout = s.cfg.Timeout
	s.server.WriteTimeout = s.cfg.Timeout

	s.logger.Info("Server has been set up")
}

func (s *ServerApp) setupLightMux() {
	s.lmux = lightmux.NewLightMux(s.server)

	mws := middlewares.NewHTTPMiddlewares(s.logger, s.cache)
	handlers := handlers.NewHTTPHandlers(s.authRepo, s.cache, s.logger)

	s.lmux.Use(mws.RecoverMiddleware)
	s.lmux.Use(mws.LoggerMiddleware)

	s.lmux.NewRoute("/ping").Handle(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("pong"))
	})
	// s.lmux.NewRoute("/panic").Handle(http.MethodGet, handlers.PanicHandler())

	authGroup := s.lmux.NewGroup("/api")
	authGroup.NewRoute("/login").Handle(http.MethodPost, handlers.LogIn())
	authGroup.NewRoute("/register").Handle(http.MethodPost, handlers.Register())
	authGroup.NewRoute("/logout", mws.JWTAuthMiddleware).Handle(http.MethodDelete, handlers.LogOut())
	authGroup.NewRoute("/refresh").Handle(http.MethodPost, handlers.RefreshTheToken())

	s.logger.Info("LightMux has been set up")
}