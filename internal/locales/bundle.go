package locales

import (
	"encoding/json"
	rootlocales "goreat/locales"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var Bundle *i18n.Bundle

func init() {
	Bundle = i18n.NewBundle(language.English)

	Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	_, err := Bundle.LoadMessageFileFS(rootlocales.FS, "en.json")
	if err != nil {
		panic(err)
	}

	_, err = Bundle.LoadMessageFileFS(rootlocales.FS, "ru.json")
	if err != nil {
		panic(err)
	}
}
