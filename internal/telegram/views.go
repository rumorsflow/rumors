package telegram

import (
	"bytes"
	"embed"
	"github.com/goccy/go-json"
	"github.com/rumorsflow/rumors/v2/internal/model"
	"github.com/rumorsflow/rumors/v2/pkg/util"
	"html/template"
	"reflect"
	"strings"
	"sync"
)

var (
	//go:embed views/*.html
	viewsFS embed.FS

	replacer = strings.NewReplacer(".", "", "-", "", " ", "")
	funcMap  = template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"join": func(data any, sep string) string {
			if isNilPtr(data) {
				return ""
			}

			switch tmp := data.(type) {
			case *[]string:
				return strings.Join(*tmp, sep)
			case []string:
				return strings.Join(tmp, sep)
			default:
				return ""
			}
		},
		"json": func(data any) any {
			if isNilPtr(data) {
				return ""
			}
			if bytes, err := json.MarshalIndent(data, "", "\t"); err == nil {
				return template.HTML(util.BytesToString(bytes))
			}
			return ""
		},
		"hashtag": func(data any) string {
			if isNilPtr(data) {
				return ""
			}

			switch tags := data.(type) {
			case []string:
				return hashtag(tags)
			case *[]string:
				return hashtag(*tags)
			default:
				return ""
			}
		},
		"domain": util.SafeDomain,
	}

	templates *template.Template
	mu        sync.Mutex
)

func init() {
	templates = template.Must(template.
		New("telegram").
		Funcs(funcMap).
		ParseFS(viewsFS, "views/*"))
}

func view(view model.View, data any) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	var out bytes.Buffer
	if err := templates.ExecuteTemplate(&out, string(view), data); err != nil {
		return "", err
	}
	return strings.Trim(out.String(), "\n"), nil
}

func isNilPtr(data any) bool {
	return reflect.TypeOf(data).Kind() == reflect.Ptr && reflect.ValueOf(data).IsNil()
}

func hashtag(data []string) string {
	var buf bytes.Buffer
	first := true
	for _, item := range data {
		if tag := replacer.Replace(item); tag != "" {
			if first {
				first = false
				_, _ = buf.WriteString("#")
			} else {
				_, _ = buf.WriteString(" #")
			}
			_, _ = buf.WriteString(tag)
		}
	}
	return buf.String()
}
