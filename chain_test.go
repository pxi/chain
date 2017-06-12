package chain

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	chain := Slice{tag("1"), tag("2"), tag("3")}
	h1 := chain.Then(serve("foo"))
	h2 := chain.Then(serve("bar"))

	rw := httptest.NewRecorder()
	h1.ServeHTTP(rw, (*http.Request)(nil))
	h2.ServeHTTP(rw, (*http.Request)(nil))

	want := "123foo123bar"
	if got := rw.Body.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func tag(s string) Func {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			io.WriteString(rw, s)
			next.ServeHTTP(rw, req)
		})
	}
}

func serve(s string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		io.WriteString(rw, s)
	})
}

func TestEmpty(t *testing.T) {
	chain := Slice{}
	h1 := chain.Then(serve("foo"))
	h2 := chain.Then(serve("bar"))

	rw := httptest.NewRecorder()
	h1.ServeHTTP(rw, (*http.Request)(nil))
	h2.ServeHTTP(rw, (*http.Request)(nil))

	want := "foobar"
	if got := rw.Body.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
