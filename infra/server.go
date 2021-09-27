package infra

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
	"gitlab.com/rendezvousPoint/incentive_calculator/config"
)

type ServerContext struct {
	Router *chi.Mux
}

func NewServerContext(server *config.Server) *ServerContext {
	router := chi.NewRouter()

	// router.Use(cors.Handler(cors.Options{
	// 	AllowedOrigins: []string{"*"},
	// 	AllowedMethods: []string{
	// 		"POST",
	// 		"GET",
	// 		"DELETE",
	// 		"PUT",
	// 		"PATCH",
	// 		"OPTIONS",
	// 	},
	// 	// AllowedHeaders
	// 	AllowedHeaders: []string{
	// 		"X-Requested-With",
	// 		"X-CSRF-Token",
	// 		"X-Auth-Token",
	// 		"Content-Type",
	// 		"Content-Length",
	// 		"processData",
	// 		"contentType",
	// 		"Origin",
	// 		"Authorization",
	// 		"Accept",
	// 		"Client-Security-Token",
	// 		"Accept-Encoding",
	// 		"timezone",
	// 		"locale",
	// 		"client_id",
	// 	},
	// 	//ExposedHeaders:   []string{"Link"},
	// 	AllowCredentials: true,
	// 	//MaxAge: 300, // Maximum value not ignored by any of major browsers
	// }))

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token", "client_id"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	router.Use(cors.Handler)
	router.Use(middleware.Timeout(time.Duration(server.ApiTimeout) * time.Second))
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	router.Use(middleware.Throttle(server.MaxServerRequests))
	router.Use(middleware.NoCache)
	router.Use(middleware.SetHeader("Content-Type", "application/json"))

	defer logrus.Info("server route setup completed")
	return &ServerContext{Router: router}
}

// StartServer starts HTTP server
func (ctx *ServerContext) StartServer(address string) {
	logrus.Info("Server starting at ", address)
	logrus.Fatal(http.ListenAndServe(address, ctx.Router))
}

func (ctx *ServerContext) Register(path string, handle func(http.ResponseWriter, *http.Request), method string) {
	ctx.Router.MethodFunc(method, path, handle)
}

func (ctx *ServerContext) Mount(path string, handler http.Handler) {
	ctx.Router.Mount(path, handler)
}
