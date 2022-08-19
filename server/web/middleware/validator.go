package middleware

import (
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
var trans ut.Translator

func init() {
	InitTrans("en")
}

func InitTrans(locale string) error {
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
		trans, _ = uni.GetTranslator(locale)

		switch locale {
		case "en":
			_ = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			_ = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			_ = enTranslations.RegisterDefaultTranslations(v, trans)
		}

	}

	return nil

}

func removeTopStruct(fields map[string]string) string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	var str string
	for _, v := range res {
		if len(str) == 0 {
			str += v
		} else {
			str += ", " + v
		}
	}

	return str
}
