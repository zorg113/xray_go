package grpc

import (
	"errors"
	"fmt"
	"net"

	"github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/app"
	"github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/logger"
	calendar_pb "github.com/zorg113/xray_go/hw12_13_14_15_calendar/pkg/calendar"
	"google.golang.org/grpc"
)

type Server struct {
	calendar_pb.UnimplementedCalendarServer
	lis    net.Listener
	l      *logger.Logger
	server *grpc.Server
	app    *app.App
}

func NewServer(log *logger.Logger, app *app.App, address string, port string) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		return nil, fmt.Errorf("start listen error: %w", err)
	}
	server := grpc.NewServer()
	srv := &Server{
		lis:    listener,
		l:      log,
		server: server,
		app:    app,
	}
	calendar_pb.RegisterCalendarServer(server, srv)
	return srv, nil
}

func (s *Server) Start() error {
	s.l.Info("start grpc server")
	if err := s.server.Serve(s.lis); err != nil {
		return fmt.Errorf("start server error: %w", err)
	}
	s.l.Info("stop grpc server")
	return nil
}

func (s *Server) Stop() error {
	if s.server == nil {
		return errors.New("grpc server in nil") //nolintlin:all
	}
	s.server.GracefulStop()
	return nil
}
