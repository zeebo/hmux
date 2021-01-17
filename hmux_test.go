package hmux

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zeebo/assert"
)

func TestShift(t *testing.T) {
	cases := []struct {
		In  string
		Dir string
		Rem string
	}{
		{"/foo/bar", "/foo", "/bar"},
		{"/foo", "/foo", ""},
		{"/", "/", ""},
		{"//", "/", "/"},
		{"//bar", "/", "/bar"},
		{"", "", ""},
	}

	for _, tc := range cases {
		dir, rem := shift(tc.In)
		assert.Equal(t, tc.Dir, dir)
		assert.Equal(t, tc.Rem, rem)
	}
}

func TestDir(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {})
	matches := func(path string, dir Dir) bool {
		rec := httptest.NewRecorder()
		dir.ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
		return rec.Code == http.StatusOK
	}

	// matches on basic functionality
	assert.That(t, matches("/foo/bar", Dir{"/foo": Dir{"/bar": ok}}))
	assert.That(t, !matches("/bar/foo", Dir{"/foo": Dir{"/bar": ok}}))

	// can use "" to assert path ends
	assert.That(t, matches("/foo/bar", Dir{"/foo": Dir{"/bar": Dir{"": ok}}}))
	assert.That(t, !matches("/foo/bar/baz", Dir{"/foo": Dir{"/bar": Dir{"": ok}}}))
	assert.That(t, !matches("/foo/bar/", Dir{"/foo": Dir{"/bar": Dir{"": ok}}}))

	// * is a wildcard match
	assert.That(t, matches("/foo/bar", Dir{"/foo": Dir{"*": ok}}))
	assert.That(t, matches("/foo/baz", Dir{"/foo": Dir{"*": ok}}))
	assert.That(t, matches("/foo", Dir{"/foo": Dir{"*": ok}}))
	assert.That(t, !matches("/foobar", Dir{"/foo": Dir{"*": ok}}))

	// * does not consume a component
	assert.That(t, matches("/foo/baz", Dir{"/foo": Dir{"*": Dir{"/baz": ok}}}))
	assert.That(t, !matches("/foo/bif", Dir{"/foo": Dir{"*": Dir{"/baz": ok}}}))

	// empty components can be matched
	assert.That(t, matches("/foo//baz", Dir{"/foo": Dir{"/": Dir{"/baz": ok}}}))

	// empty key only matches empty url
	assert.That(t, !matches("/", Dir{"": ok}))
}

func TestMethod(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {})
	matches := func(method string, m Method) bool {
		rec := httptest.NewRecorder()
		m.ServeHTTP(rec, httptest.NewRequest(method, "/", nil))
		return rec.Code == http.StatusOK
	}

	// matches on basic functionality
	assert.That(t, matches("POST", Method{"POST": ok}))
	assert.That(t, !matches("PUT", Method{"POST": ok}))
}

func TestArg(t *testing.T) {
	arg1 := Arg("1")
	arg2 := Arg("2")

	// check that argument shifts and captures the path component
	arg1.Capture(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "foo", arg1.Value(req.Context()))
		assert.Equal(t, "/bar", req.URL.Path)
	})).ServeHTTP(nil, httptest.NewRequest("", "/foo/bar", nil))

	// check that empty argument works
	Dir{"/foo": arg1.Capture(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "", arg1.Value(req.Context()))
		assert.Equal(t, "/bar", req.URL.Path)
	}))}.ServeHTTP(nil, httptest.NewRequest("", "/foo//bar", nil))

	// check that no argument is a 404
	rec := httptest.NewRecorder()
	Dir{"/foo": arg1.Capture(nil)}.ServeHTTP(rec, httptest.NewRequest("", "/foo", nil))
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// check double arguments don't get confused
	arg1.Capture(arg2.Capture(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "foo", arg1.Value(req.Context()))
		assert.Equal(t, "bar", arg2.Value(req.Context()))
		assert.Equal(t, "", req.URL.Path)
	}))).ServeHTTP(nil, httptest.NewRequest("", "/foo/bar", nil))

	// check that argument capture is value oriented
	arg1.Capture(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "foo", Arg("1").Value(req.Context()))
		assert.Equal(t, "/bar", req.URL.Path)
	})).ServeHTTP(nil, httptest.NewRequest("", "/foo/bar", nil))
}
