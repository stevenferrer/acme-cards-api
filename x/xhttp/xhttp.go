package xhttp

import (
	"log/slog"
	"net/http"
)

// Handler is the same as http.Handler except ServeHTTP may return an error.
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request) error
}

// HandlerFunc is a convenience type like http.HandlerFunc.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// ServeHTTP implements the Handler interface.
func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return h(w, r)
}

func WrapXHTTP(h Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := h.ServeHTTP(w, r)
		if err != nil {
			handleError(w, r, err)
		}
	})
}

func handleError(w http.ResponseWriter, _ *http.Request, err error) {
	// TODO: Inject logger into context
	slog.Error("http handler error", "err", err)

	// TODO: Send appropriate error status based on error
	w.WriteHeader(http.StatusInternalServerError)
}
