// SPDX-FileCopyrightText: © 2019 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirector_ServeHTTP(t *testing.T) {
	r := &redirector{"src.example.com/x", "https://example.com/git", "git"}

	t.Run("GoVisit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "https://src.example.com/foo?go-get=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		resp := w.Result()
		if http.StatusOK != resp.StatusCode {
			t.Errorf("expected %d, got %d", http.StatusFound, resp.StatusCode)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		expected := `<!doctype html>
<meta name="go-import" content="src.example.com/foo git https://example.com/git/src.example.com">
<title>go-import-redirect</title>
`
		if string(body) != expected {
			t.Errorf("expected\n---\n%s\n---\ngot\n---\n%s\n---", expected, string(body))
		}
		if hdr := resp.Header.Get("Location"); hdr != "" {
			t.Error("expected empty Location header")
		}
	})

	t.Run("UserVisit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "https://src.example.com/foo", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		resp := w.Result()
		if http.StatusFound != resp.StatusCode {
			t.Errorf("expected %d, got %d", http.StatusFound, resp.StatusCode)
		}
		if resp.ContentLength > 0 {
			t.Error("expected empty body")
		}
		if hdr := resp.Header.Get("Location"); hdr != "https://pkg.go.dev/src.example.com/foo" {
			t.Errorf("expected %q, got %q", "https://pkg.go.dev/src.example.com/foo", hdr)
		}
	})
}

func TestRedirector_getRepo(t *testing.T) {
	r := &redirector{"src.example.com/x/", "https://example.com/git/", "git"}
	for _, tc := range []struct{ pkg, expected string }{
		{"src.example.com/x/foo", "https://example.com/git/foo"},
		{"src.example.com/x/foo/bar", "https://example.com/git/foo"},
	} {
		if actual := r.getRepo(tc.pkg); actual != tc.expected {
			t.Errorf("expected %q, got %q", tc.expected, actual)
		}
	}

	r = &redirector{"src.example.com/x", "https://example.com/git", "git"}
	for _, tc := range []struct{ pkg, expected string }{
		{"src.example.com/x/foo", "https://example.com/git/foo"},
		{"src.example.com/x/foo/bar", "https://example.com/git/foo"},
	} {
		if actual := r.getRepo(tc.pkg); actual != tc.expected {
			t.Errorf("expected %q, got %q", tc.expected, actual)
		}
	}
}
