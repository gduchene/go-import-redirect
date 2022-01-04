// SPDX-FileCopyrightText: © 2019 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

package main

import (
	"log"
	"net/http"
	"path"
	"strings"
	"text/template"
)

var body = template.Must(template.New("").Parse(`<!doctype html>
<meta name="go-import" content="{{.Package}} {{.VCS}} {{.Repository}}">
<title>go-import-redirect</title>
`))

type bodyData struct{ Package, Repository, VCS string }

type redirector struct{ from, to, vcs string }

var _ http.Handler = &redirector{}

func (h *redirector) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pkg := path.Join(req.Host, req.URL.Path)
	if req.URL.Query().Get("go-get") != "1" {
		w.Header().Set("Location", "https://pkg.go.dev/"+pkg)
		w.WriteHeader(http.StatusFound)
		return
	}
	dest := h.getRepo(pkg)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := body.Execute(w, bodyData{pkg, dest, h.vcs}); err != nil {
		log.Println(err)
	}
}

func (h *redirector) getRepo(pkg string) string {
	from := strings.TrimRight(h.from, "/")
	to := strings.TrimRight(h.to, "/")
	path := strings.TrimLeft(strings.TrimPrefix(pkg, from), "/")
	return to + "/" + strings.Split(path, "/")[0]
}
