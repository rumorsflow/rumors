package task

import (
	"context"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/otiai10/opengraph/v2"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
	"github.com/spf13/cast"
	"golang.org/x/net/html"
	"net/http"
	"strings"
	"sync"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/111.0"

var regexMap sync.Map

func addRegex(r *string) error {
	if r != nil && *r != "" {
		if _, ok := regexMap.Load(r); !ok {
			re, err := regexp2.Compile(*r, regexp2.IgnoreCase&regexp2.RE2)
			if err != nil {
				return err
			}
			regexMap.Store(r, re)
		}
	}
	return nil
}

func matchByLoc(r *string, str string) bool {
	if re, ok := regexMap.Load(r); ok {
		ok, _ := re.(*regexp2.Regexp).MatchString(str)
		return ok
	}
	return true
}

func searchByLoc(r *string, str string) string {
	if re, ok := regexMap.Load(r); ok {
		if m, _ := re.(*regexp2.Regexp).FindStringMatch(str); m != nil {
			return m.String()
		}
	}
	return ""
}

func openGraphFetch(ctx context.Context, url string) (*opengraph.OpenGraph, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if !strings.HasPrefix(res.Header.Get("Content-Type"), "text/html") {
		return nil, errs.New("content type must be text/html")
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("open graph error due to request %s with response status code %d", url, res.StatusCode)
	}

	og := opengraph.New(url)
	og.Intent.TrustedTags = []string{opengraph.HTMLMetaTag, opengraph.HTMLTitleTag, opengraph.HTMLLinkTag}
	node, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	if err = walk(og, node); err != nil {
		return nil, err
	}

	return og, nil
}

func walk(og *opengraph.OpenGraph, node *html.Node) error {
	if node.Type == html.ElementNode {
		switch {
		case node.Data == opengraph.HTMLMetaTag && trust(og, node.Data):
			return MetaTag(node).Contribute(og)
		case node.Data == opengraph.HTMLTitleTag && trust(og, node.Data):
			return opengraph.TitleTag(node).Contribute(og)
		case node.Data == opengraph.HTMLLinkTag && trust(og, node.Data):
			return opengraph.LinkTag(node).Contribute(og)
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		walk(og, child)
	}

	return nil
}

func trust(og *opengraph.OpenGraph, tagName string) bool {
	for _, name := range og.Intent.TrustedTags {
		if name == tagName {
			return true
		}
	}
	return false
}

func toMedia(og *opengraph.OpenGraph) []entity.Media {
	media := make([]entity.Media, 0, len(og.Image)+len(og.Video)+len(og.Audio))
	for _, i := range og.Image {
		media = append(media, entity.Media{URL: i.URL, Type: entity.ImageType, Meta: map[string]any{
			"width":  i.Width,
			"height": i.Height,
			"alt":    i.Alt,
		}})
	}
	for _, i := range og.Video {
		media = append(media, entity.Media{URL: i.URL, Type: entity.VideoType, Meta: map[string]any{
			"width":    i.Width,
			"height":   i.Height,
			"duration": i.Duration,
		}})
	}
	for _, i := range og.Audio {
		media = append(media, entity.Media{URL: i.URL, Type: entity.AudioType})
	}
	return media
}

func pagination(args string) (i uint64, s uint64, search string) {
	s = 10

	if args == "" {
		return
	}

	a := strings.Fields(args)
	if len(a) > 0 {
		i = cast.ToUint64(a[0])
	}

	if len(a) > 1 {
		if s = cast.ToUint64(a[1]); s == 0 {
			s = 10
		}
		if s > 20 {
			s = 20
		}
	}

	if len(a) > 2 {
		search = strings.TrimSpace(strings.Join(a[2:], " "))
	}
	return
}

func contains(data []string, el string) bool {
	for _, item := range data {
		if strings.EqualFold(item, el) {
			return true
		}
	}
	return false
}

type Meta struct {
	*opengraph.Meta
}

func MetaTag(node *html.Node) *Meta {
	meta := new(opengraph.Meta)
	for _, attr := range node.Attr {
		switch attr.Key {
		case "property":
			meta.Property = attr.Val
		case "content":
			meta.Content = attr.Val
		case "name":
			meta.Name = attr.Val
		}
	}
	return &Meta{Meta: meta}
}

func (meta *Meta) Contribute(og *opengraph.OpenGraph) (err error) {
	switch meta.Name {
	case "parsely-title":
		og.Title = meta.Content
	case "parsely-link":
		og.URL = meta.Content
	case "parsely-image-url":
		if len(og.Image) == 0 || og.Image[len(og.Image)-1].URL != meta.Content {
			og.Image = append(og.Image, opengraph.Image{URL: meta.Content})
		}
	default:
		return meta.Meta.Contribute(og)
	}
	return nil
}
