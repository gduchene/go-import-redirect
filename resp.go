// SPDX-FileCopyrightText: © 2019 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

package main

import (
	_ "embed"
	"log"
	"net/http"
	"path"
	"regexp"
	"text/template"
)

var (
	body = template.Must(template.New("").Parse(tmpl))

	//go:embed resp.html
	tmpl string
)

type bodyData struct{ Package, Repository, VCS string }

type redirector struct {
	re   *regexp.Regexp
	repl string
	vcs  string
}

func (h *redirector) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pkg := path.Join(req.Host, req.URL.Path)
	if !h.re.MatchString(pkg) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if req.URL.Query().Get("go-get") != "1" {
		w.Header().Set("Location", "https://pkg.go.dev/"+pkg)
		w.WriteHeader(http.StatusFound)
		return
	}

	dest := h.re.ReplaceAllString(pkg, h.repl)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := body.Execute(w, bodyData{pkg, dest, h.vcs}); err != nil {
		log.Println(err)
	}
}
