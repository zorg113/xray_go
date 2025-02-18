package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/logger"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Server struct {
	Address string
	server  *http.Server
	logger  logger.Logger
}

// type Logger interface {
// 	Info(msg string)
// 	Error(msg string)
// 	Warn(msg string)
// 	Debug(msg string) // TODO
// }

type Application interface { // TODO
}

func NewServer(host, port string, log logger.Logger, app Application) *Server {
	return &Server{
		Address: net.JoinHostPort(host, port),
		logger:  log,
	}
}

func (s *Server) Start(ctx context.Context) error {
	router := mux.NewRouter()
	router.HandleFunc("/hi", s.helloWorld).Methods("GET")
	router.Use(s.loggingMiddleware)
	s.server = &http.Server{
		Addr:         s.Address,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "start server error")
	}

	select {
	case <-ctx.Done():
		return nil
	}
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return errors.New("server is nil")
	}
	if err := s.server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "stop server error")
	}
	return nil
}

func (s *Server) helloWorld(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}

// TODO
