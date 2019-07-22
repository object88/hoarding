package http

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/object88/hoarding"
)

func Test_Http_ViaContext(t *testing.T) {
	tcs := []struct {
		name string
	}{
		{
			name: "WithContext",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

		})
	}

	var b bytes.Buffer

	h := NewHandler(func(h *Handler) *Handler {
		h.hc = func() *hoarding.Hoard {
			return hoarding.NewHoarder(&b, func() bool {
				return true
			})
		}
		return h
	})

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h, _ := r.Context().Value(HoarderContextKey).(*hoarding.Hoard)
		h.Msgf(hoarding.Info, "msg1")
		w.Write([]byte("OK"))
	})

	h0 := &http.Server{Addr: ":9999", Handler: h.ServeHTTP(final)}

	go func() {
		h0.ListenAndServe()
	}()

	resp, err := http.Get("http://:9999")
	if err != nil {
		t.Fatalf("Received unexpected error: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	h0.Shutdown(context.Background())

	actual := b.String()
	if actual != "msg1" {
		t.Errorf("Did not recieve the anticipated message: '%s'", actual)
	}
}

type TestCompositionHandler struct {
	h *Handler
}

func (tch *TestCompositionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tch.h.Msgf(hoarding.Info, "msg1")
	w.Write([]byte("OK"))
}

func Test_Http_ViaComposition(t *testing.T) {
	var b bytes.Buffer

	h := NewHandler(func(h *Handler) *Handler {
		h.hc = func() *hoarding.Hoard {
			return hoarding.NewHoarder(&b, func() bool {
				return true
			})
		}
		return h
	})

	h0 := &http.Server{Addr: ":9999", Handler: h.ServeHTTP(&TestCompositionHandler{h: h})}

	go func() {
		h0.ListenAndServe()
	}()

	resp, err := http.Get("http://:9999")
	if err != nil {
		t.Fatalf("Received unexpected error: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	h0.Shutdown(context.Background())

	actual := b.String()
	if actual != "msg1" {
		t.Errorf("Did not recieve the anticipated message: '%s'", actual)
	}
}
