package views

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"strings"
	"sync"
)

type ViewNS int

const (
	TelegramNS ViewNS = iota + 1
)

var (
	//go:embed telegram/*.html
	telegramFS embed.FS

	templates map[ViewNS]*template.Template
	mu        sync.Mutex
)

var funcMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	},
	"join": func(data []string, sep string) string {
		return strings.Join(data, sep)
	},
}

func init() {
	templates = make(map[ViewNS]*template.Template)
	templates[TelegramNS] = template.Must(template.New("telegram").Funcs(funcMap).ParseFS(telegramFS, "telegram/*"))
}

func Render(w io.Writer, ns ViewNS, template string, data any) error {
	mu.Lock()
	defer mu.Unlock()

	return templates[ns].ExecuteTemplate(w, template, data)
}

func View(ns ViewNS, template string, data any) (string, error) {
	var out bytes.Buffer
	if err := Render(&out, ns, template, data); err != nil {
		return "", err
	}
	return out.String(), nil
}
