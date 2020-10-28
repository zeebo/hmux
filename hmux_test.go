package hmux

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
		dir.ServeHTTP(rec, &http.Request{URL: &url.URL{Path: path}})
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

	// * does not consume a component
	assert.That(t, matches("/foo/baz", Dir{"/foo": Dir{"*": Dir{"/baz": ok}}}))
	assert.That(t, !matches("/foo/bif", Dir{"/foo": Dir{"*": Dir{"/baz": ok}}}))

	// empty components can be matched
	assert.That(t, matches("/foo//baz", Dir{"/foo": Dir{"/": Dir{"/baz": ok}}}))
}

func TestMethod(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {})
	matches := func(method string, m Method) bool {
		rec := httptest.NewRecorder()
		m.ServeHTTP(rec, &http.Request{Method: method})
		return rec.Code == http.StatusOK
	}

	// matches on basic functionality
	assert.That(t, matches("POST", Method{"POST": ok}))
	assert.That(t, !matches("PUT", Method{"POST": ok}))
}

func TestArg(t *testing.T) {
	arg := new(Arg)
	arg2 := new(Arg)

	// check that argument shifts and captures the path component
	arg.Capture(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		assert.That(t, arg.Exists(req.Context()))
		assert.Equal(t, "/foo", arg.Value(req.Context()))

		assert.Equal(t, "/bar", req.URL.Path)
	})).ServeHTTP(nil, httptest.NewRequest("", "/foo/bar", nil))

	// check that no argument is a 404
	rec := httptest.NewRecorder()
	arg.Capture(nil).ServeHTTP(rec, httptest.NewRequest("", "/", nil))
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// check double arguments don't get confused
	arg.Capture(arg2.Capture(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		assert.That(t, arg.Exists(req.Context()))
		assert.Equal(t, "/foo", arg.Value(req.Context()))

		assert.That(t, arg2.Exists(req.Context()))
		assert.Equal(t, "/bar", arg2.Value(req.Context()))

		assert.Equal(t, "", req.URL.Path)
	}))).ServeHTTP(nil, httptest.NewRequest("", "/foo/bar", nil))
}
