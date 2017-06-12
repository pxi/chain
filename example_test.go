package chain_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pxi/chain"
)

// timeout returns a middleware function that will timeout after d.
func timeout(d time.Duration) chain.Func {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, d, "timed out")
	}
}

// timed returns a middleware function that prints handling times to w.
func timed(w io.Writer) chain.Func {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, req *http.Request) {
			start := time.Now()
			next.ServeHTTP(rw, req)
			fmt.Fprintf(w, "%s %s took %s\n", req.Method, req.URL, time.Since(start))
		}
		return http.HandlerFunc(fn)
	}
}

func Example() {
	// Construct a default middleware chain.
	common := chain.Slice{timeout(2 * time.Second)}

	// Add additional middleware for some paths.
	logged := append(common, timed(os.Stderr))

	handler := (http.Handler)(nil)

	mux := http.NewServeMux()
	mux.Handle("/foo", logged.Then(handler))
	mux.Handle("/", common.Then(handler))

	// When the mux is getting served, all requests are timed out after 2
	// seconds. Additionally requests to /foo are getting their handling times
	// printed to stdout as well.
}
