package rest

import (
	"context"
	"net/http"

	"github.com/dexlabsio/garlic/logging"
	chi "github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var servers map[string]*Server

type Server struct {
	Name string
	root chi.Router
}

func NewServer(name string) *Server {
	router := chi.NewRouter()

	return &Server{
		Name: name,
		root: router,
	}
}

func (s *Server) Router() chi.Router {
	return s.root
}

// Listen starts an HTTP server on the specified bind address and listens for incoming requests.
// It runs in a separate goroutine and returns a channel to report errors. The server will log
// the bind address and continue running until the context is canceled, at which point it will
// send the context's error to the error channel. Note that http.ListenAndServe does not stop
// immediately upon context cancellation; it stops only on an error or when the process exits.
func (s *Server) Listen(ctx context.Context, bind string) <-chan error {
	l := logging.Global()
	errCh := make(chan error, 1)

	// This version immediately stops reporting when the context is canceled,
	// but note that http.ListenAndServe itself doesn't stop immediately—it
	// stops only on an error or when the process exits.
	go func() {
		defer close(errCh)

		go func() {
			l.Info("Listening.", zap.String("bind", bind), zap.String("server", s.Name))
			err := http.ListenAndServe(bind, s.Router())
			errCh <- err
		}()

		<-ctx.Done()
		errCh <- ctx.Err()
	}()

	return errCh
}

// GetServer implements a multiton of servers
func GetServer(name string) *Server {
	if servers == nil {
		servers = make(map[string]*Server)
	}

	if srv, exists := servers[name]; exists {
		return srv
	}

	srv := NewServer(name)
	servers[name] = srv
	return srv
}
