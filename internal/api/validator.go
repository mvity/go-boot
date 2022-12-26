/*
 * Copyright © 2021 - 2022 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package api

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh_Hans_SG"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/zh"
	"github.com/mvity/go-boot/internal/app"
	"github.com/mvity/go-box/x"
	"reflect"
	"strings"
)

// 初始化自定义验证器
func initValidator() {

	if valid, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 初始化gin中文错误消息
		chinese := zh_Hans_SG.New()
		uti := ut.New(chinese, chinese)
		app.Trans, _ = uti.GetTranslator("zh")
		_ = zh.RegisterDefaultTranslations(valid, app.Trans)

		// 绑定自定义tag `label`
		valid.RegisterTagNameFunc(func(fld reflect.StructField) string {
			tag := fld.Tag.Get("label")
			if tag == "" {
				tag = "[" + fld.Tag.Get("json") + "]"
			} else {
				tag += "[" + fld.Tag.Get("json") + "]"
			}
			return tag
		})

		// 注册手机号验证器
		_ = valid.RegisterValidation("mobile", func(fl validator.FieldLevel) bool {
			val := strings.TrimSpace(fl.Field().String())
			if val == "" {
				return true
			}
			return x.RegexpChinaMobile.MatchString(val)
		})
		_ = valid.RegisterTranslation("mobile", app.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0}"+"格式无效", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})

	}
}
