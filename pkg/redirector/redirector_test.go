package redirector_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"go.awhk.org/core"
	"go.awhk.org/go-import-redirect/pkg/redirector"
)

func TestRedirector_ServeHTTP(s *testing.T) {
	t := core.T{T: s}

	meta := regexp.MustCompile(`<meta name="go-import" content="(.+?)">`)
	redir := &redirector.Redirector{
		Template: redirector.DefaultTemplate,
		Transformers: []redirector.Transformer{
			{
				Pattern:     &redirector.Pattern{regexp.MustCompile("go.example.com/x/baz(?:/.+)?")},
				Replacement: "https://git.example.net/elsewhere/baz",
				VCS:         "git",
			},
			{
				Pattern:     &redirector.Pattern{regexp.MustCompile("go.example.com/x/([^/]+)(?:/.+)?")},
				Replacement: "https://git.example.net/y/$1",
				VCS:         "git",
			},
			{
				Pattern:     &redirector.Pattern{regexp.MustCompile("go.example.com/x/qux(?:/.+)?")},
				Replacement: "https://git.example.net/elsewhere/qux",
				VCS:         "git",
			},
		},
	}

	for _, tc := range []struct {
		name   string
		method string
		url    string

		expGetStatusCode   int
		expGetGoImport     string
		expVisitStatusCode int
		expVisitLocation   string
	}{
		{
			name:   "Match",
			method: http.MethodGet,
			url:    "https://go.example.com/x/foo",

			expGetStatusCode:   http.StatusOK,
			expGetGoImport:     "go.example.com/x/foo git https://git.example.net/y/foo",
			expVisitStatusCode: http.StatusFound,
			expVisitLocation:   "https://pkg.go.dev/go.example.com/x/foo",
		},
		{
			name:   "MatchDirectory",
			method: http.MethodGet,
			url:    "https://go.example.com/x/foo/bar",

			expGetStatusCode:   http.StatusOK,
			expGetGoImport:     "go.example.com/x/foo/bar git https://git.example.net/y/foo",
			expVisitStatusCode: http.StatusFound,
			expVisitLocation:   "https://pkg.go.dev/go.example.com/x/foo/bar",
		},
		{
			name:   "MatchIgnored",
			method: http.MethodGet,
			url:    "https://go.example.com/x/qux",

			expGetStatusCode:   http.StatusOK,
			expGetGoImport:     "go.example.com/x/qux git https://git.example.net/y/qux",
			expVisitStatusCode: http.StatusFound,
			expVisitLocation:   "https://pkg.go.dev/go.example.com/x/qux",
		},
		{
			name:   "MatchSpecific",
			method: http.MethodGet,
			url:    "https://go.example.com/x/baz",

			expGetStatusCode:   http.StatusOK,
			expGetGoImport:     "go.example.com/x/baz git https://git.example.net/elsewhere/baz",
			expVisitStatusCode: http.StatusFound,
			expVisitLocation:   "https://pkg.go.dev/go.example.com/x/baz",
		},
		{
			name:   "BadMethod",
			method: http.MethodPost,
			url:    "https://go.example.com/x/baz",

			expGetStatusCode:   http.StatusMethodNotAllowed,
			expVisitStatusCode: http.StatusMethodNotAllowed,
		},
	} {
		t.Run("Get"+tc.name, func(t *core.T) {
			var (
				req = httptest.NewRequest(tc.method, tc.url+"?go-get=1", nil)
				w   = httptest.NewRecorder()
			)
			redir.ServeHTTP(w, req)

			resp := w.Result()
			t.AssertEqual(tc.expGetStatusCode, resp.StatusCode)
			t.AssertEqual("", resp.Header.Get("Location"))

			match := meta.FindSubmatch(core.Must(io.ReadAll(resp.Body)))
			if tc.expGetGoImport == "" {
				t.AssertEqual(0, len(match))
				return
			}
			if t.AssertEqual(2, len(match)) {
				t.AssertEqual(tc.expGetGoImport, string(match[1]))
			}
		})
		t.Run("Visit"+tc.name, func(t *core.T) {
			var (
				req = httptest.NewRequest(tc.method, tc.url, nil)
				w   = httptest.NewRecorder()
			)
			redir.ServeHTTP(w, req)

			resp := w.Result()
			t.AssertEqual(tc.expVisitStatusCode, resp.StatusCode)
			t.AssertEqual(tc.expVisitLocation, resp.Header.Get("Location"))
		})
	}
}
