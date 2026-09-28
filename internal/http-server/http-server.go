package httpserver

import (
	"context"
	"log/slog"
	"net/http"

	genericjwtservice "github.com/ayayaakasvin/generic-jwt-service"
	"github.com/ayayaakasvin/oneflick-ticket/internal/config"
	"github.com/ayayaakasvin/oneflick-ticket/internal/domain"
	"github.com/ayayaakasvin/oneflick-ticket/internal/http-server/handlers"
	"github.com/ayayaakasvin/oneflick-ticket/internal/http-server/middlewares"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/ayayaakasvin/lightmux"
)

type ServerApp struct {
	server  *http.Server
	httpcfg *config.HTTPServerConfig
	corscfg *config.CorsConfig
	lmux    *lightmux.LightMux

	ar   domain.UserRepository
	er   domain.EventRepository
	cc   domain.Cache
	lfs  domain.FS
	rlm  *domain.RateLimiter
	jwtM *genericjwtservice.JWTService
	smtp domain.SMTP

	logger *slog.Logger
}

func NewServerApp(
	httpconfig *config.HTTPServerConfig,
	corsconfig *config.CorsConfig,
	logger *slog.Logger,
	authRepo domain.UserRepository,
	eventRepo domain.EventRepository,
	cache domain.Cache,
	lfs domain.FS,
	rlm *domain.RateLimiter,
	jwtM *genericjwtservice.JWTService,
	smtp domain.SMTP,
) *ServerApp {
	return &ServerApp{
		httpcfg: httpconfig,
		corscfg: corsconfig,
		logger:  logger,
		ar:      authRepo,
		er:      eventRepo,
		cc:      cache,
		lfs:     lfs,
		rlm:     rlm,
		jwtM:    jwtM,
		smtp:    smtp,
	}
}

func (s *ServerApp) Start(ctx context.Context) error {
	s.setupServer()
	s.setupLightMux()

	return func() error {
		s.logger.Info("Server has been started", "port", s.httpcfg.Address)

		return s.lmux.Run(ctx)
	}()
}

// func (s *ServerApp) StartTLS(ctx context.Context) error {
// 	s.setupServer()
// 	s.setupLightMux()

// 	return func() error {
// 		s.logger.Info("Server has been started", "port", s.httpcfg.Address)

// 		return s.lmux.RunTLS(ctx, s.tlsCfg.CertFile, s.tlsCfg.KeyFile)
// 	}()
// }

func (s *ServerApp) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// setuping server by pointer, so we dont have to return any value
func (s *ServerApp) setupServer() {
	if s.server == nil {
		s.logger.Warn("Server is nil, creating a new server pointer")
		s.server = &http.Server{}
	}

	s.server.Addr = s.httpcfg.Address
	s.server.IdleTimeout = s.httpcfg.IdleTimeout
	s.server.ReadTimeout = s.httpcfg.Timeout
	s.server.WriteTimeout = s.httpcfg.Timeout

	s.logger.Info("Server has been set up")
}

func (s *ServerApp) setupLightMux() {
	s.lmux = lightmux.NewLightMux(s.server)

	mws := middlewares.NewHTTPMiddlewares(s.logger, s.cc, s.jwtM, s.corscfg, s.ar, s.rlm)
	hndlrs := handlers.NewHTTPHandlers(s.ar, s.er, s.cc, s.lfs, s.smtp, s.jwtM, s.logger)

	s.lmux.Use(mws.CaptureMiddleware, mws.RecoverMiddleware, mws.MetricsMiddleware, mws.LoggerMiddleware, mws.CORSMiddleware)
	s.lmux.NewRoute("/metrics").Handle(http.MethodGet, promhttp.Handler().ServeHTTP) // TODO: dont forget to wrap!
	s.lmux.NewRoute("/ping").Handle(http.MethodGet, hndlrs.PingHandler())

	apiRoute := s.lmux.NewGroup("/api")

	authGroup := apiRoute.ContinueGroup("/auth")
	authGroup.NewRoute("/register").Handle(http.MethodPost, hndlrs.Register())                    // /api/auth/register POST
	authGroup.NewRoute("/login").Handle(http.MethodPost, hndlrs.Login())                          // /api/auth/login POST
	authGroup.NewRoute("/logout", mws.JWTAuthMiddleware).Handle(http.MethodPost, hndlrs.Logout()) // /api/auth/logout POST
	authGroup.NewRoute("/refresh").Handle(http.MethodPost, hndlrs.Refresh())                      // /api/auth/refresh POST

	// Event group with JWT middleware
	eventGroup := apiRoute.ContinueGroup("/event", mws.JWTAuthMiddleware)
	eventGroup.NewRoute("/all").Handle(http.MethodGet, hndlrs.GetAllEvents())               // all events
	eventGroup.NewRoute("/category").Handle(http.MethodGet, hndlrs.GetEventsByCategoryID()) // events by category_id
	// get top 10
	eventGroup.NewRoute("/top-ten").Handle(http.MethodGet, hndlrs.GetTop10Events())

	// event CRUD
	eventRoute := eventGroup.NewRoute("")
	eventRoute.Handle(http.MethodPost, hndlrs.SaveEvent())                                                        // save
	eventGroup.NewRoute("/update/upload").Handle(http.MethodPost, hndlrs.UpdateEventImageURLByUploading())        // image upload
	eventGroup.NewRoute("/update/image").Handle(http.MethodPost, hndlrs.UpdateEventImageURLUsingExternalSource()) // external image link
	eventRoute.Handle(http.MethodGet, hndlrs.GetEventByUUID())                                                    // event by event_uuid
	eventRoute.Handle(http.MethodDelete, hndlrs.DeleteEventByUUID())                                              // delete by event_uuid

	// Category Route GET method
	apiRoute.ContinueGroup("/category", mws.JWTAuthMiddleware).NewRoute("").Handle(http.MethodGet, hndlrs.GetAllCategories()) // all category ids and names
	apiRoute.ContinueGroup("/images/", mws.JWTAuthMiddleware).NewRoute("").Handle(http.MethodGet, hndlrs.ServeImages())       // images saved in main server, served by custom FS

	// Ticket CRUD
	ticketGroup := apiRoute.ContinueGroup("/ticket")
	ticketRoute := ticketGroup.NewRoute("")
	ticketRoute.Handle(http.MethodGet, hndlrs.GetTicket())
	ticketRoute.Handle(http.MethodPost, hndlrs.InsertTicketAfterwards())
	ticketRoute.Handle(http.MethodDelete, hndlrs.DeleteTicket())

	s.lmux.Mux().HandleFunc("/docs/", httpSwagger.WrapHandler)

	s.lmux.Mux().HandleFunc("/", hndlrs.NotFound())

	s.logger.Info("LightMux has been set up")
}
