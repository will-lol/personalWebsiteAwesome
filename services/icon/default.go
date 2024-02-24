package icon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/will-lol/personalWebsiteAwesome/lib/pointerify"
)

type versions struct {
	Svg  []string `json:"svg"`
	Font []string `json:"font"`
}

type alias struct {
	Original string `json:"base"`
	Alias    string `json:"alias"`
}

type icon struct {
	Name     string   `json:"name"`
	AltNames []string `json:"altnames"`
	Tags     []string `json:"tags"`
	Versions versions `json:"versions"`
	Color    string   `json:"color"`
	Aliases  []alias `json:"aliases"`
}

type iconFinder struct {
	Icons map[string]icon
}

type IconFinder interface {
	Find(name string) (url *string, err error)
}

func NewIconFinder(iconJsonEndpoint string) (IconFinder, error) {
	resp, err := http.Get(iconJsonEndpoint)
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var i []icon
	err = json.Unmarshal(bytes, &i)
	if err != nil {
		return nil, err
	}

	m := make(map[string]icon, len(i))
	for _, icon := range i {
		m[icon.Name] = icon
	}

	return &iconFinder{
		Icons: m,
	}, nil
}

func (i iconFinder) Find(name string) (url *string, err error) {
	const original = "original"
	const wordmark = "original-wordmark"

	icon, ok := i.Icons[name]
	if ok == false {
		return nil, errors.New("Icon does not exist")
	}

	if slices.Contains(icon.Versions.Svg, original) {
		return pointerify.Pointer(fmt.Sprintf("https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/%s/%s-%s.svg", icon.Name, icon.Name, original)), err
	}

	if slices.Contains(icon.Versions.Svg, wordmark) {
		return pointerify.Pointer(fmt.Sprintf("https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/%s/%s-%s.svg", icon.Name, icon.Name, wordmark)), err
	}

	return nil, errors.New("Desired version not found")
}
