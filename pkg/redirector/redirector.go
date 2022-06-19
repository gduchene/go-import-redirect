package redirector

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"
	"text/template"

	"go.awhk.org/core"
)

var DefaultTemplate = template.Must(template.ParseFS(fs, "*.tpl"))

var (
	filter = core.FilterHTTPMethod(http.MethodGet)

	//go:embed *.tpl
	fs embed.FS
)

type Redirector struct {
	Template     *template.Template
	Transformers []Transformer
}

func (h *Redirector) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if filter(w, req) {
		return
	}

	pkg := path.Join(req.Host, req.URL.Path)
	for _, t := range h.Transformers {
		if !t.Pattern.MatchString(pkg) {
			continue
		}

		if req.URL.Query().Get("go-get") != "1" {
			w.Header().Set("Location", "https://pkg.go.dev/"+pkg)
			w.WriteHeader(http.StatusFound)
			return
		}

		repo := t.Pattern.ReplaceAllString(pkg, t.Replacement)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if err := h.Template.Execute(w, TemplateData{pkg, repo, t.VCS}); err != nil {
			log.Println("Failed to execute template:", err)
		}
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

type Pattern struct{ *regexp.Regexp }

func (pp *Pattern) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err = json.Unmarshal(data, &s); err != nil {
		return
	}
	pp.Regexp, err = regexp.Compile(strings.ReplaceAll(s, `\\`, `\`))
	return
}

type TemplateData struct{ Package, Repository, VCS string }

type Transformer struct {
	Pattern     *Pattern `json:"pattern"`
	Replacement string   `json:"replacement"`
	VCS         string   `json:"vcs"`
}
