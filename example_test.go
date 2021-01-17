package hmux_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/zeebo/hmux"
)

func Example() {
	handler := http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		fmt.Printf("method: %q\n", req.Method)
		fmt.Printf("path: %q\n", req.URL.Path)
		fmt.Printf("name: %q\n", hmux.Arg("name").Value(req.Context()))
	})

	resources := hmux.Dir{
		"/foo": hmux.Dir{
			"*": hmux.Arg("name").Capture(
				hmux.Method{
					"GET":  handler,
					"POST": handler,
				},
			),
		},
	}

	resources.ServeHTTP(nil, httptest.NewRequest("POST", "/foo/bar", nil))
	fmt.Println("---")
	resources.ServeHTTP(nil, httptest.NewRequest("GET", "/foo/baz/bif", nil))

	//output:
	// method: "POST"
	// path: ""
	// name: "bar"
	// ---
	// method: "GET"
	// path: "/bif"
	// name: "baz"
}

func ExampleMethod() {
	resources := hmux.Method{
		"POST": http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
			fmt.Printf("method: %q\n", req.Method)
			fmt.Printf("path: %q\n", req.URL.Path)
		}),
	}

	resources.ServeHTTP(nil, httptest.NewRequest("POST", "/foo", nil))

	//output:
	// method: "POST"
	// path: "/foo"
}

func ExampleDir() {
	resources := hmux.Dir{
		"/foo": http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
			fmt.Printf("path: %q\n", req.URL.Path)
		}),
	}

	resources.ServeHTTP(nil, httptest.NewRequest("GET", "/foo/bar", nil))

	//output:
	// path: "/bar"
}

func ExampleArg() {
	resources := hmux.Arg("name").Capture(
		http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
			fmt.Printf("arg: %q\n", hmux.Arg("name").Value(req.Context()))
			fmt.Printf("path: %q\n", req.URL.Path)
		}),
	)

	resources.ServeHTTP(nil, httptest.NewRequest("GET", "/foo/bar", nil))

	//output:
	// arg: "foo"
	// path: "/bar"
}
