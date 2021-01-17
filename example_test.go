package hmux_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/zeebo/hmux"
)

func ExampleMethod() {
	resources := hmux.Method{
		"POST": http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
			fmt.Println("method:", req.Method)
			fmt.Println("path:", req.URL.Path)
		}),
	}

	resources.ServeHTTP(nil, httptest.NewRequest("POST", "/foo", nil))

	//output:
	// method: POST
	// path: /foo
}

func ExampleDir() {
	resources := hmux.Dir{
		"/foo": http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
			fmt.Println("mux: /foo")
			fmt.Println("path:", req.URL.Path)
		}),
	}

	resources.ServeHTTP(nil, httptest.NewRequest("GET", "/foo/bar", nil))

	//output:
	// mux: /foo
	// path: /bar
}

func ExampleArg() {
	resources := hmux.Arg("name").Capture(
		http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
			fmt.Println("arg:", hmux.Arg("name").Value(req.Context()))
			fmt.Println("path:", req.URL.Path)
		}),
	)

	resources.ServeHTTP(nil, httptest.NewRequest("GET", "/foo/bar", nil))

	//output:
	// arg: foo
	// path: /bar
}
