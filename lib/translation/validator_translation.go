package translation

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

func InitTrans(locale string) (ut.Translator, error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			if name == "" {
				name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
				if name == "-" {
					return ""
				}
			}
			return name
		})
		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT, zhT, enT)
		trans, _ := uni.GetTranslator(locale)

		switch locale {
		case "en":
			_ = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			_ = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			_ = enTranslations.RegisterDefaultTranslations(v, trans)
		}

		return trans, nil
	}

	return nil, errors.New("init translation error")

}
