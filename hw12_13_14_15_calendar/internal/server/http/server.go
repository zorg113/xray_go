package internalhttp

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/app"
	"github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/logger"
	"github.com/zorg113/xray_go/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	Address string
	server  *http.Server
	logger  logger.Logger
	app     *app.App
}

func NewServer(host, port string, log logger.Logger, app *app.App) *Server {
	return &Server{
		Address: net.JoinHostPort(host, port),
		logger:  log,
		app:     app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	router := mux.NewRouter()
	router.HandleFunc("/hello", s.helloWorld).Methods("GET")
	router.HandleFunc("/create", s.createEvent).Methods("POST")
	router.HandleFunc("/update", s.updateEvent).Methods("POST")
	router.HandleFunc("/get", s.getEvents).Methods("GET")
	router.HandleFunc("/delete", s.deleteEvent).Methods("POST")
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

	<-ctx.Done()
	return nil
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

func (s *Server) helloWorld(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	var event storage.Event
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("cannot decode to struct" + err.Error())
		_, err = w.Write([]byte("cannot decode to struct" + err.Error()))
		if err != nil {
			s.logger.Error("cannot write to reply" + err.Error())
			return
		}
		return
	}
	if err := s.app.CreateEvent(event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("cannot create event" + err.Error()))
		s.logger.Error("cannot create event" + err.Error())
		if err != nil {
			s.logger.Error("cannot write to reply" + err.Error())
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Event created")); err != nil {
		s.logger.Error("cannot write to reply" + err.Error())
		return
	}
}

func (s *Server) updateEvent(w http.ResponseWriter, r *http.Request) {
	var event storage.Event
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("cannot decode to struct" + err.Error())
		if _, err := w.Write([]byte("cannot decode to struct")); err != nil {
			s.logger.Error("cannot write to reply" + err.Error())
			return
		}
		return
	}
	if err := s.app.UpdateEvent(event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("cannot update event" + err.Error())
		if _, err := w.Write([]byte("cannot update event")); err != nil {
			s.logger.Error("cannot write to reply" + err.Error())
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Event created")); err != nil {
		s.logger.Error("cannot write to reply" + err.Error())
		return
	}
}

func (s *Server) getEvents(w http.ResponseWriter, r *http.Request) {
	var event storage.Event
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("cannot decode to struct" + err.Error())
		if _, err := w.Write([]byte("cannot decode to struct")); err != nil {
			s.logger.Error("cannot write to reply" + err.Error())
			return
		}
		return
	}
	events, err := s.app.GetEvents(event.StartDate, event.EndDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("cannot get events" + err.Error())
		if _, err := w.Write([]byte("cannot get events")); err != nil {
			s.logger.Error("cannot write to reply" + err.Error())
			return
		}
		return
	}

	if len(events) == 0 {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("events not found")); err != nil {
			s.logger.Error("cannot write to reply" + err.Error())
			return
		}
		return
	}

	result, err := json.Marshal(events)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("cannot marshal events" + err.Error())
		if _, err := w.Write([]byte("cannot marshal events")); err != nil {
			s.logger.Error("cannot write to reply" + err.Error())
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("cannot write to reply" + err.Error())
		return
	}
}

func (s *Server) deleteEvent(w http.ResponseWriter, r *http.Request) {
	var event storage.Event
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("cannot decode to struct" + err.Error())
		if _, err := w.Write([]byte("cannot decode to struct")); err != nil {
			s.logger.Error("cannot write to reply" + err.Error())
			return
		}
		return
	}
	if err := s.app.DeleteEvent(event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("cannot delete event" + err.Error())
		if _, err := w.Write([]byte("cannot delete event")); err != nil {
			s.logger.Error("cannot write to reply" + err.Error())
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Event deleted")); err != nil {
		s.logger.Error("cannot write to reply" + err.Error())
		return
	}
}
