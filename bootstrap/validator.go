package bootstrap

import (
	"reflect"
	"shalabing-gin/utils"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitializeValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册自定义验证器
		// 通用校验--------------------------------------------
		_ = v.RegisterValidation("mobile", utils.ValidateMobile)
		// 后台管理员校验--------------------------------------------
		_ = v.RegisterValidation("manager_username", utils.ValidateManagerUsername)
		_ = v.RegisterValidation("manager_password", utils.ValidateManagerPassword)
		_ = v.RegisterValidation("exist_manager", utils.ValidateManagerByIdExist)
		_ = v.RegisterValidation("exist_manager_username", utils.ValidateManagerUsernameExist)
		_ = v.RegisterValidation("exist_role", utils.ValidateRoleExist)
		_ = v.RegisterValidation("exist_permission", utils.ValidatePermissionExist)
		_ = v.RegisterValidation("permission_slice", utils.ValidatePermissionSlice)
		_ = v.RegisterValidation("dept_slice", utils.ValidateDeptSlice)
		_ = v.RegisterValidation("exist_dept", utils.ValidateDeptExist)
		_ = v.RegisterValidation("exist_post", utils.ValidatePostExist)
		// 用户端校验--------------------------------------------

		// 注册自定义 json tag 函数
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}
