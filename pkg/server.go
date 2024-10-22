package pkg

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/zapr"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server is the server
type Server struct {
	HttpPort         int
	BindAddress      string
	engine           *gin.Engine
	hServer          *http.Server
	shutdownComplete chan bool
}

// NewServer creates a new server
func NewServer(httpPort int) (*Server, error) {
	s := &Server{
		HttpPort:    httpPort,
		BindAddress: "0.0.0.0",
	}

	if err := s.createGinEngine(); err != nil {
		return nil, err
	}
	return s, nil
}

// Run starts the http server
func (s *Server) Run() error {
	s.shutdownComplete = make(chan bool, 1)
	trapInterrupt(s)

	log := zapr.NewLogger(zap.L())

	if s.HttpPort <= 0 {
		return errors.New("HTTP port must be a positive integer")
	}
	address := fmt.Sprintf("%s:%d", s.BindAddress, s.HttpPort)
	log.Info("Starting http server", "address", address)

	s.hServer = &http.Server{
		Handler: s.engine,
	}

	lis, err := net.Listen("tcp", address)

	if err != nil {
		return errors.Wrapf(err, "Could not start listener")
	}
	go func() {
		if err := s.hServer.Serve(lis); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error(err, "There was an error with the http server")
			}
		}
	}()

	// Wait for the shutdown to complete
	// We use a channel to signal when the shutdown method has completed and then return.
	// This is necessary because shutdown() is running in a different go function from hServer.Serve. So if we just
	// relied on hServer.Serve to return and then returned from Run we might still be in the middle of calling shutdown.
	// That's because shutdown calls hServer.Shutdown which causes hserver.Serve to return.
	<-s.shutdownComplete
	return nil
}

// createGinEngine sets up the gin engine which is a router
func (s *Server) createGinEngine() error {
	log := zapr.NewLogger(zap.L())
	log.Info("Setting up the server")

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, world!")
	})
	router.GET("/healthz", s.healthCheck)

	s.engine = router
	return nil
}

func (s *Server) shutdown() {
	log := zapr.NewLogger(zap.L())
	log.Info("Shutting down the Rube server")

	if s.hServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		if err := s.hServer.Shutdown(ctx); err != nil {
			log.Error(err, "Error shutting down http server")
		}
		log.Info("HTTP Server shutdown complete")
	}
	log.Info("Shutdown complete")
	s.shutdownComplete <- true
}

// trapInterrupt shutdowns the server if the appropriate signals are sent
func trapInterrupt(s *Server) {
	log := zapr.NewLogger(zap.L())
	sigs := make(chan os.Signal, 10)
	// Note SIGSTOP and SIGTERM can't be caught
	// We can trap SIGINT which is what ctl-z sends to interrupt the process
	// to interrupt the process
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		msg := <-sigs
		log.Info("Received signal", "signal", msg)
		s.shutdown()
	}()
}

func (s *Server) healthCheck(ctx *gin.Context) {
	// TODO(jeremy): We should return the version
	d := gin.H{
		"server": "rube",
		"status": "healthy",
	}
	code := http.StatusOK
	ctx.JSON(code, d)
}
