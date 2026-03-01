package http_mw

import (
	"goreat/internal/locales"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const LocalizerKey = "Localizer"

func LocalizerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptLang := c.GetHeader("Accept-Language")
		localizer := i18n.NewLocalizer(locales.Bundle, acceptLang)
		c.Set(LocalizerKey, localizer)

		c.Next()
	}
}

func GetLocalizer(c *gin.Context) *i18n.Localizer {
	loc, exists := c.Get(LocalizerKey)
	if !exists {
		return nil
	}
	return loc.(*i18n.Localizer)
}
