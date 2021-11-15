package messages

import (
	"embed"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Embed a directory
//go:embed translates
var embedDirTranslates embed.FS

const _dir = "translates"

func NewBundle() (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	fs, err := embedDirTranslates.ReadDir(_dir)
	if err != nil {
		return bundle, err
	}

	for _, f := range fs {
		p := f.Name()
		log.Println(p)
		if strings.HasSuffix(p, ".toml") && strings.HasPrefix(p, "active") {
			data, err := embedDirTranslates.ReadFile(path.Join(_dir, p))
			if err != nil {
				return bundle, err
			}

			if _, err := bundle.ParseMessageFileBytes(data, p); err != nil {
				return bundle, err
			}
		}
	}

	return bundle, nil
}

func MustNewBundle() *i18n.Bundle {
	b, err := NewBundle()
	if err != nil {
		panic(fmt.Sprintf("load localizer bundle failed: %v", err))
	}
	return b
}
