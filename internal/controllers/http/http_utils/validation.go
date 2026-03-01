package http_utils

import (
	"errors"
	"goreat/internal/controllers/http/http_mw"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

			if name == "-" {
				return ""
			}
			return name
		})
	}
}

func msgForTag(fe validator.FieldError, localizer *i18n.Localizer) string {
	var templateData map[string]interface{}

	switch fe.Tag() {
	case "oneof":
		templateData = map[string]interface{}{"Param": fe.Param()}
	case "min":
		templateData = map[string]interface{}{"Param": fe.Param()}
	}

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    fe.Tag(),
		TemplateData: templateData,
	})

	if err != nil {
		return fe.Error()
	}
	return msg
}

func formatValidationErrors(err error, localizer *i18n.Localizer) map[string]string {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		out := make(map[string]string)
		for _, fe := range ve {
			out[fe.Field()] = msgForTag(fe, localizer)
		}
		return out
	}

	invalidReqMsg, _ := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "InvalidRequest",
	})
	if invalidReqMsg == "" {
		invalidReqMsg = "Invalid data format"
	}

	return map[string]string{"request": invalidReqMsg}
}

func ShouldBindJSON(q any, ctx *gin.Context) map[string]string {
	if err := ctx.ShouldBindJSON(q); err != nil {
		return formatValidationErrors(err, http_mw.GetLocalizer(ctx))
	}
	return nil
}
