package httpc

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTimeoutOptionPassesThroughBeforeDeadline(t *testing.T) {
	opt := &TimeoutOption{Timeout: 200 * time.Millisecond}
	handler := opt.Next(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test", "ok")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("done"))
	})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	if recorder.Code != http.StatusCreated {
		t.Fatalf("unexpected status: got %d want %d", recorder.Code, http.StatusCreated)
	}
	if recorder.Header().Get("X-Test") != "ok" {
		t.Fatalf("unexpected header: got %q want %q", recorder.Header().Get("X-Test"), "ok")
	}
	if recorder.Body.String() != "done" {
		t.Fatalf("unexpected body: got %q want %q", recorder.Body.String(), "done")
	}
}

func TestTimeoutOptionReturns408AndCancelsHandler(t *testing.T) {
	opt := &TimeoutOption{Timeout: 20 * time.Millisecond}

	done := make(chan error, 1)
	handler := opt.Next(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
		_, err := w.Write([]byte("too late"))
		done <- err
	})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	if recorder.Code != http.StatusRequestTimeout {
		t.Fatalf("unexpected status: got %d want %d", recorder.Code, http.StatusRequestTimeout)
	}
	if recorder.Body.Len() != 0 {
		t.Fatalf("unexpected body: got %q want empty", recorder.Body.String())
	}

	select {
	case err := <-done:
		if !errors.Is(err, ErrHandlerTimeout) {
			t.Fatalf("unexpected write error: got %v want %v", err, ErrHandlerTimeout)
		}
	case <-time.After(time.Second):
		t.Fatal("handler did not observe timeout cancellation")
	}
}

func TestTimeoutOptionPropagatesPanicBeforeTimeout(t *testing.T) {
	opt := &TimeoutOption{Timeout: time.Second}
	handler := opt.Next(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	})

	defer func() {
		if p := recover(); p == nil {
			t.Fatal("expected panic to propagate")
		}
	}()

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
}
