package validate

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var trans ut.Translator

func InitTranslator() {
	zhTrans := zh.New()
	uni := ut.New(zhTrans, zhTrans)
	trans, _ = uni.GetTranslator("zh")
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		registerTranslations(v, trans)
	}
	binding.EnableDecoderUseNumber = true
}
func registerTranslations(v *validator.Validate, trans ut.Translator) {
	_ = v.RegisterTranslation("required", trans,
		func(ut ut.Translator) error {
			return ut.Add("required", "{0}不能为空", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required", fe.Field())
			return t
		})
	_ = v.RegisterTranslation("min", trans,
		func(ut ut.Translator) error {
			return ut.Add("min", "{0}长度不足", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("min", fe.Field(), fe.Param())
			return t
		})
	_ = v.RegisterTranslation("max", trans,
		func(ut ut.Translator) error {
			return ut.Add("max", "{0}长度超出限制", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("max", fe.Field(), fe.Param())
			return t
		})
	_ = v.RegisterTranslation("email", trans,
		func(ut ut.Translator) error {
			return ut.Add("email", "{0}邮箱格式不正确", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("email", fe.Field())
			return t
		})
}

// TranslateError 将验证错误翻译成中文
func TranslateError(err error) string {
	if err == nil {
		return ""
	}

	var ve validator.ValidationErrors
	if errs, ok := err.(validator.ValidationErrors); ok {
		ve = errs
	} else {
		return err.Error()
	}
	for _, e := range ve {
		return e.Translate(trans)
	}
	return ""
}
