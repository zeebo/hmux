# package hmux

`import "github.com/zeebo/hmux"`

<p>
  <a href="https://pkg.go.dev/github.com/zeebo/hmux"><img src="https://img.shields.io/badge/doc-reference-007d9b?logo=go&style=flat-square" alt="go.dev" /></a>
  <a href="https://goreportcard.com/report/github.com/zeebo/hmux"><img src="https://goreportcard.com/badge/github.com/zeebo/hmux?style=flat-square" alt="Go Report Card" /></a>
  <a href="https://sourcegraph.com/github.com/zeebo/hmux?badge"><img src="https://sourcegraph.com/github.com/zeebo/hmux/-/badge.svg?style=flat-square" alt="SourceGraph" /></a>
</p>



## Usage

#### type Arg

```go
type Arg struct {
}
```

Arg captures path components and attaches them to the request context. It always
captures a non-empty component.

#### func (*Arg) Capture

```go
func (a *Arg) Capture(h http.Handler) http.Handler
```
Capture consumes a path component and stores it in the request context so that
it can be retreived with Value.

#### func (*Arg) Value

```go
func (a *Arg) Value(ctx context.Context) string
```
Value returns the value associated with the Arg on the context.

#### type Dir

```go
type Dir map[string]http.Handler
```

Dir pulls off path components from the front of the path and dispatches. It
attempts to dispatch to "*" without consuming a path component if nothing
matches.

#### func (Dir) ServeHTTP

```go
func (d Dir) ServeHTTP(w http.ResponseWriter, req *http.Request)
```
ServeHTTP implements the http.Handler interface.

#### type Method

```go
type Method map[string]http.Handler
```

Method checks the request method and dispatches.

#### func (Method) ServeHTTP

```go
func (m Method) ServeHTTP(w http.ResponseWriter, req *http.Request)
```
ServeHTTP implements the http.Handler interface.
