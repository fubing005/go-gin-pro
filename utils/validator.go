package utils

import (
	"regexp"
	"shalabing-gin/app/models"
	"shalabing-gin/global"

	"github.com/go-playground/validator/v10"
)

// 通用校验--------------------------------------------------------------

// ValidateMobile 校验手机号
func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	ok, _ := regexp.MatchString(`^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}$`, mobile)
	if !ok {
		return false
	}
	return true
}

// 后台管理员校验-----------------------------------------------------------

// 校验后台管理员用户名
// 举一个例子[^[a-zA-Z0-9]{5,8}$]
func ValidateManagerUsername(fl validator.FieldLevel) bool {
	userName := fl.Field().String()
	ok, _ := regexp.MatchString(`^[a-z]{5,10}$`, userName)
	if !ok {
		return false
	}
	return true
}

// 校验后台管理员密码
func ValidateManagerPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	ok, _ := regexp.MatchString(`^[0-9]{5,10}$`, password)
	if !ok {
		return false
	}
	return true
}

// 验证后台管理员用户名是否存在
func ValidateManagerUsernameExist(fl validator.FieldLevel) bool {
	// id := fl.Field().Uint()
	userName := fl.Field().String()

	//判断管理员是否存在
	// if id > 0 { //编辑
	// 	manager := models.Manager{}
	// 	global.App.DB.Where("id != ? AND username = ?", id, userName).Find(&manager)
	// 	if manager.ID.ID > 0 {
	// 		return false
	// 	}
	// } else { //添加
	manager := models.Manager{}
	global.App.DB.Where("username = ?", userName).Find(&manager)
	//如果id > 0 ,说明数据库存在该管理员，则返回false
	if manager.ID.ID > 0 {
		return false
	}
	// }

	return true
}

// 验证后台管理员用户名是否存在
func ValidateManagerByIdExist(fl validator.FieldLevel) bool {
	ID := fl.Field().Uint()

	//判断管理员是否存在
	manager := models.Manager{}
	global.App.DB.Where("id = ? AND is_super = 0", ID).Find(&manager)
	//如果id > 0 ,说明数据库存在该管理员，则返回false
	if manager.ID.ID > 0 {
		return true
	}

	return false
}

// 验证是否存在后台角色
func ValidateRoleExist(fl validator.FieldLevel) bool {
	roleId := fl.Field().Uint()

	if roleId == 1 {
		return false
	}

	//判断角色是否存在
	role := models.Role{}
	global.App.DB.Where("id = ?", roleId).Find(&role)
	//如果id > 0 ,说明数据库存在该角色则返回true
	if role.ID.ID > 0 {
		return true
	}

	return false
}

// 验证是否存在后台部门
func ValidateDeptExist(fl validator.FieldLevel) bool {
	deptId := fl.Field().Uint()

	dept := models.Dept{}
	global.App.DB.Where("id = ?", deptId).Find(&dept)
	if dept.ID.ID > 0 {
		return true
	}

	return false
}

// 验证是否存在后台岗位
func ValidatePostExist(fl validator.FieldLevel) bool {
	postId := fl.Field().Uint()

	post := models.Post{}
	global.App.DB.Where("id = ?", postId).Find(&post)
	if post.ID.ID > 0 {
		return true
	}

	return false
}

// 验证是否存在相关权限,如果存在返回true则可以编辑
func ValidatePermissionExist(fl validator.FieldLevel) bool {
	permissionId := fl.Field().Uint()

	//判断角色是否存在
	permission := models.Permission{}
	global.App.DB.Where("id = ?", permissionId).Find(&permission)
	//如果id > 0 ,说明数据库存在该权限则返回true
	if permission.ID.ID > 0 {
		return true
	}

	return false
}

// 定义一个函数来验证切片中的元素是否都不是零值
func ValidatePermissionSlice(fl validator.FieldLevel) bool {
	// 将FieldLevel中的值转换为切片
	permissionSlice, ok := fl.Field().Interface().([]uint)
	if ok {
		// 遍历切片中的每个元素
		for _, v := range permissionSlice {
			// 如果有元素是零值，返回false
			if v == 0 {
				return false
			}
			//判断权限是否存在
			permission := models.Permission{}
			global.App.DB.Where("id = ?", v).Find(&permission)
			if permission.ID.ID == 0 {
				return false
			}
		}

		// 所有元素都不是零值，返回true
		return true
	}
	// 如果转换失败，返回false
	return false
}

// 定义一个函数来验证切片中的元素是否都不是零值
func ValidateDeptSlice(fl validator.FieldLevel) bool {
	// 将FieldLevel中的值转换为切片
	deptSlice, ok := fl.Field().Interface().([]uint)
	if ok {
		// 遍历切片中的每个元素
		for _, v := range deptSlice {
			// 如果有元素是零值，返回false
			if v == 0 {
				return false
			}
			//判断部门是否存在
			dept := models.Dept{}
			global.App.DB.Where("id = ?", v).Find(&dept)
			if dept.ID.ID == 0 {
				return false
			}
		}

		// 所有元素都不是零值，返回true
		return true
	}
	// 如果转换失败，返回false
	return false
}

//用户端校验------------------------------------------------------
