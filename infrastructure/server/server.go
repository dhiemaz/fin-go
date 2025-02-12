package server

import (
	"bitbucket.org/rctiplus/almasbub"
	"fmt"
	"github.com/dhiemaz/fin-go/config"
	"github.com/dhiemaz/fin-go/infrastructure/logger"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	HTTP_PORT                    = "8081" // HTTP_PORT default value
	HTTP_MAX_CONN_PER_IP         = 5000   // HTTP_MAX_CONN_PER_IP server default value
	HTTP_MAX_REQUEST_PER_CONN    = 100    // HTTP_MAX_REQUEST_PER_CONN server default value
	HTTP_MAX_CONCURRENCY         = 100000 // HTTP_MAX_CONCURRENCY server default value
	HTTP_MAX_KEEP_ALIVE_DURATION = 500    // HTTP_MAX_KEEP_ALIVE_DURATION server default value
)

type server struct {
	http   *fasthttp.Server
	router *router.Router
}

// StartServer : start HTTP Server
func Start(router *router.Router) {
	var wg sync.WaitGroup

	// create new instance server
	server := newServer(router)

	// run server
	server.Run(&wg)
	wg.Wait()
}

// newServer : creates instance server
func newServer(router *router.Router) *server {
	// create handler from router
	handler := router.Handler

	// define fasthttp server
	return &server{
		http:   newHTTPServer(handler),
		router: router,
	}
}

// newHTTPServer creates a new HTTP Server
func newHTTPServer(h fasthttp.RequestHandler) *fasthttp.Server {

	// set maxConnsPerIP from config, if failed then set default value
	maxConnsPerIP := almasbub.ToInt(config.GetConfig().HTTPMaxConnPerIP)
	if maxConnsPerIP == 0 {
		maxConnsPerIP = HTTP_MAX_CONN_PER_IP // set default value
	}

	// set maxRequestsPerConn from config, if failed then set default value
	maxRequestsPerConn := almasbub.ToInt(config.GetConfig().HTTPMaxRequestPerConn)
	if maxRequestsPerConn == 0 {
		maxRequestsPerConn = HTTP_MAX_REQUEST_PER_CONN // set default value
	}

	// set concurrency from config, if failed then set default value
	concurrency := almasbub.ToInt(config.GetConfig().HTTPMaxConcurrency)
	if concurrency == 0 {
		concurrency = HTTP_MAX_CONCURRENCY // The maximum number of concurrent requests http may process
	}

	// set maxKeepaliveDuration from config, if failed then set default value
	maxKeepalive := almasbub.ToDuration(config.GetConfig().HTTPMaxKeepAlive)
	if maxKeepalive == 0 {
		maxKeepalive = HTTP_MAX_KEEP_ALIVE_DURATION // set dafault value
	}

	return &fasthttp.Server{
		Handler:              h,
		MaxConnsPerIP:        maxConnsPerIP,
		MaxRequestsPerConn:   maxRequestsPerConn,
		MaxKeepaliveDuration: maxKeepalive * time.Millisecond,
		MaxRequestBodySize:   1024 * 1024 * 1024 * 4,
		Concurrency:          concurrency,
		ReduceMemoryUsage:    true,
	}
}

// Run starts the HTTP server and performs a graceful shutdown
func (s *server) Run(wg *sync.WaitGroup) {
	port := almasbub.ToString(config.GetConfig().Port)
	if port == "" {
		port = HTTP_PORT // set http port default value (8081)
	}

	logger.WithFields(logger.Fields{"component": "server", "action": "run http server", "port": port}).
		Infof("run http server.")

	// create a fast listener ;)
	ln, err := reuseport.Listen("tcp4", fmt.Sprintf(":%s", port))
	if err != nil {
		logger.WithFields(logger.Fields{"component": "server", "action": "run http server", "port": port}).
			Fatalf("reuseport listener failed, error : %v", err)
	}

	// create a graceful shutdown listener
	duration := 2 * time.Second
	graceful := NewGracefulListener(ln, duration)

	// Error handling
	listenErr := make(chan error, 1)

	// Run server
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.WithFields(logger.Fields{"component": "server", "action": "run http server", "port": port}).
			Infof("starting tusd uploader server, listening on port %v", port)
		listenErr <- s.http.Serve(graceful)
	}()

	// SIGINT/SIGTERM handling
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	// Handle channels/graceful shutdown
	for {
		select {
		// If server.ListenAndServe() cannot start due to errors such
		// as "port in use" it will return an error.
		case err := <-listenErr:
			if err != nil {
				logger.WithFields(logger.Fields{"component": "server", "action": "run http server", "port": port}).
					Infof("server failed to start, error : %v", err)
			}
			os.Exit(0)

		// handle termination signal
		case <-osSignals:
			logger.WithFields(logger.Fields{"component": "server", "action": "stopping http server", "port": port}).
				Infof("shutdown signal received")

			// Servers in the process of shutting down should disable KeepAlives
			s.http.DisableKeepalive = true

			// Attempt the graceful shutdown by closing the listener
			// and completing all inflight requests.
			if err := graceful.Close(); err != nil {
				logger.WithFields(logger.Fields{"component": "server", "action": "stopping http server", "port": port}).
					Errorf("server gracefully stop with error : %v", err)
			} else {
				logger.WithFields(logger.Fields{"component": "server", "action": "stopping http server", "port": port}).
					Infof("server gracefully stop")
			}
		}
	}
}
