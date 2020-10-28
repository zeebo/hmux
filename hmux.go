package hmux

import (
	"context"
	"net/http"
	"strings"
)

// Method checks the request method and dispatches.
type Method map[string]http.Handler

// ServeHTTP implements the http.Handler interface.
func (m Method) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if h, ok := m[req.Method]; ok {
		h.ServeHTTP(w, req)
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// shift removes the first path component from a string.
func shift(path string) (dir, left string) {
	if len(path) == 0 || path[0] != '/' {
		return path, ""
	} else if split := strings.IndexByte(path[1:], '/') + 1; split == 0 {
		return path, ""
	} else {
		return path[:split], path[split:]
	}
}

// Dir pulls off path components from the front of the path and dispatches.
// It attempts to dispatch to "*" without consuming a path component if nothing matches.
type Dir map[string]http.Handler

// ServeHTTP implements the http.Handler interface.
func (d Dir) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	dir, rem := shift(req.URL.Path)
	if h, ok := d[dir]; ok {
		req.URL.Path = rem
		h.ServeHTTP(w, req)
	} else if h, ok := d["*"]; ok {
		h.ServeHTTP(w, req)
	} else {
		http.Error(w, "not found", http.StatusNotFound)
	}
}

// Arg captures path components and attaches them to the request context
type Arg struct {
	byte // non-zero sized so that pointers are distinct
}

// Value returns the value associated with the Arg on the context. It returns
// the zero value if it is not set.
func (a *Arg) Value(ctx context.Context) string { return getArguments(ctx)[a] }

// Exists returns true if the Arg has a value on the context.
func (a *Arg) Exists(ctx context.Context) bool {
	_, ok := getArguments(ctx)[a]
	return ok
}

// Capture consumes a path component and stores it in the request context so that it
// can be retreived with Value.
func (a *Arg) Capture(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		dir, rem := shift(req.URL.Path)
		if req.URL.Path != "" && req.URL.Path != "/" {
			req.URL.Path = rem
			h.ServeHTTP(w, req.WithContext(addArgument(req.Context(), a, dir)))
		} else {
			http.Error(w, "not found", http.StatusNotFound)
		}
	})
}

//
// we store a map[*Arg]string and reuse it on the context to avoid O(N) lookup behavior
//

type argumentsKey struct{}

func getArguments(ctx context.Context) map[*Arg]string {
	args, _ := ctx.Value(argumentsKey{}).(map[*Arg]string)
	return args
}

func addArgument(ctx context.Context, a *Arg, val string) context.Context {
	if args := getArguments(ctx); args != nil {
		args[a] = val
		return ctx
	}
	return context.WithValue(ctx, argumentsKey{}, map[*Arg]string{a: val})
}
