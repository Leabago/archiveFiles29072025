package httpserver

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const defaultPort = ":8080"
const defaultTimout = 30 * time.Second

type Server struct {
	Server *http.Server
	notify chan error

	port            string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutDownTimeout time.Duration
}

func New(port string) *Server {
	s := &Server{
		port:            defaultPort,
		readTimeout:     defaultTimout,
		writeTimeout:    defaultTimout,
		shutDownTimeout: defaultTimout,
	}

	s.Server = &http.Server{
		Addr:         s.port,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
	}

	return s
}

// старт
func (s *Server) Start(l *zap.Logger) {
	go func() {
		l.Info("start server", zap.String("port", s.port))
		s.notify <- s.Server.ListenAndServe()
		close(s.notify)
	}()
}

// нотифицировать ошибкой если сервер не запустился
func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown(l *zap.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	l.Debug("http shutdown")
	return s.Server.Shutdown(ctx)
}
