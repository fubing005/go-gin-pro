package services_admin

import (
	"encoding/json"
	"errors"
	"shalabing-gin/app/models"
	"shalabing-gin/core/trans"
	"shalabing-gin/global"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type mainService struct{}

var MainService = new(mainService)

func (mainService *mainService) Index(c *gin.Context) (data map[string]interface{}, err error) {
	//1.获取用户信息
	userinfoStr, ok := c.Keys["userinfo"].(string)
	if ok {
		//1.获取用户信息
		var userinfoStruct models.Manager
		//把获取到的用户信息转换结构体
		json.Unmarshal([]byte(userinfoStr), &userinfoStruct)
		//获取所有权限列表
		permissionList := []models.Permission{}
		global.App.DB.Where("module_id = ?", 0).Preload("PermissionItem", func(db *gorm.DB) *gorm.DB {
			return db.Order("permission.sort ASC")
		}).Order("sort ASC").Find(&permissionList)

		//获取当前角色拥有的权限,并把权限id放在一个map对象中
		rolePermission := []models.RolePermission{}
		global.App.DB.Where("role_id = ?", userinfoStruct.RoleId).Find(&rolePermission)

		dept_id, ok := c.Keys["dept_id"].(uint)
		if ok {
			//判断当前用户所属部门，根据部门查询部门所属角色，将查询到的角色id放在rolePermission
			roleIds := []uint{}
			global.App.DB.Model(&models.RoleDept{}).Where("dept_id = ?", dept_id).Select("role_id").Find(&roleIds)
			roleDeptPermission := []models.RolePermission{}
			global.App.DB.Where("role_id in (?)", roleIds).Find(&roleDeptPermission)
			rolePermission = append(rolePermission, roleDeptPermission...)
		}

		rolePermissionMap := make(map[uint]uint)
		for _, v := range rolePermission {
			rolePermissionMap[v.PermissionID] = v.PermissionID
		}

		//循环遍历所有权限数据,判断当前权限的id是否在角色权限的map对象中,如果是的话给当前数据加入checked属性
		for i := 0; i < len(permissionList); i++ { //循环权限列表
			if userinfoStruct.IsSuper == 1 { //超级管理员
				permissionList[i].Checked = true
				for j := 0; j < len(permissionList[i].PermissionItem); j++ { // 判断当前权限的子栏位是否在角色权限的map中
					permissionList[i].PermissionItem[j].Checked = true
					// for jj := 0; jj < len(permissionList[i].PermissionItem[j].PermissionItem); jj++ { // 判断当前权限的子栏位是否在角色权限的map中
					// 	permissionList[i].PermissionItem[j].PermissionItem[jj].Checked = true
					// }
				}
			} else { //不是超级管理员
				if _, ok := rolePermissionMap[permissionList[i].ID.ID]; ok { // 判断当前权限是否在角色权限的map对象中
					permissionList[i].Checked = true
				}
				for j := 0; j < len(permissionList[i].PermissionItem); j++ { // 判断当前权限的子栏位是否在角色权限的map中
					if _, ok := rolePermissionMap[permissionList[i].PermissionItem[j].ID.ID]; ok { // 判断当前权限是否在角色权限的map对象中
						permissionList[i].PermissionItem[j].Checked = true
					}
					// for jj := 0; jj < len(permissionList[i].PermissionItem[j].PermissionItem); jj++ { // 判断当前权限的子栏位是否在角色权限的map中
					// 	if _, ok := rolePermissionMap[permissionList[i].PermissionItem[j].PermissionItem[jj].ID.ID]; ok { // 判断当前权限是否在角色权限的map对象中
					// 		permissionList[i].PermissionItem[j].PermissionItem[jj].Checked = true
					// 	}
					// }
				}
			}
		}

		data = map[string]interface{}{
			"username":        userinfoStruct.Username,
			"is_super":        userinfoStruct.IsSuper,
			"permission_list": permissionList,
		}
		return data, nil
	} else {
		return nil, errors.New(trans.Trans("common.登录授权失效"))
	}
}
