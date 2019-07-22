package http

import (
	"context"
	"net/http"

	"github.com/object88/hoarding"
)

// ContextKey is a key into a context
type ContextKey int

const (
	// HoarderContextKey is the key into the context which accesses the Hoard
	HoarderContextKey ContextKey = 0
)

// Handler is HTTP middleware for the hoarding logger
type Handler struct {
	hoarding.Hoard
	// Log hoarding.Level
	hc HoardCreator
}

type HoardCreator func() *hoarding.Hoard

type Option func(*Handler) *Handler

// NewHandler constructs a Handler
func NewHandler(opts ...Option) *Handler {
	h := &Handler{
		hc: func() *hoarding.Hoard {
			return hoarding.NewHoarder(nil, func() bool {
				return true
			})
		},
	}

	for _, o := range opts {
		h = o(h)
	}

	return h
}

func (h *Handler) ServeHTTP(h0 http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Hoard = *h.hc()

		r = r.WithContext(context.WithValue(r.Context(), HoarderContextKey, &h.Hoard))
		h0.ServeHTTP(w, r)

		h.Flush()
	})
}
