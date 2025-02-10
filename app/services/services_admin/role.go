package services_admin

import (
	"errors"
	"shalabing-gin/app/common/request"
	"shalabing-gin/app/common/request/request_admin"
	"shalabing-gin/app/models"
	"shalabing-gin/core/trans"
	"shalabing-gin/global"
	"shalabing-gin/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type roleService struct{}

var RoleService = new(roleService)

var roleCount int64

func (roleService *roleService) Common() (data map[string]interface{}) {
	data = make(map[string]interface{}, 2)
	status := []Status{{Code: 1, Name: "启用"}, {Code: 2, Name: "禁用"}}
	data["status"] = status
	return
}

// 分页获取角色
func (roleService *roleService) Index(form request.PageQuery, c *gin.Context) (data map[string]interface{}) {
	data = make(map[string]interface{}, 2)

	// managerId, _ := c.Keys["id"].(uint)
	// isSuper := c.Keys["is_super"].(int)

	// data = make(map[string]interface{}, 2)
	roleList := []models.Role{}
	// if isSuper == 1 { // 超级管理员拥有所有角色
	global.App.DB.Offset((form.Page - 1) * form.PageSize).Limit(form.PageSize).Find(&roleList)
	data["list"] = roleList

	global.App.DB.Model(&models.Role{}).Count(&roleCount)
	data["count"] = roleCount
	// } else {
	// 	// 如果manager dept_id 的 parent_id = 0, 则获取该部门下的所有用户创建的角色
	// 	deptId := c.Keys["dept_id"].(uint)
	// 	parentId := 0
	// 	global.App.DB.Model(&models.Dept{}).Where("id = ?", deptId).Select("parent_id").Find(&parentId)

	// 	if parentId == 0 {
	// 		// 获取当前用户的子部门ID
	// 		deptSubIds := []uint{}
	// 		global.App.DB.Model(&models.Dept{}).Where("parent_id IN (?)", deptId).Select("id").Find(&deptSubIds)

	// 		// 获取当前用户的子部门的所有用户ID
	// 		managerSubIds := []uint{}
	// 		global.App.DB.Model(&models.Manager{}).Where("dept_id IN (?)", deptSubIds).Select("id").Find(&managerSubIds)

	// 		// 获取当前用户创建的角色，以及当前用户的子用户创建的角色
	// 		global.App.DB.Where("create_by = ?", managerId).Or("create_by IN (?)", managerSubIds).Offset((form.Page - 1) * form.PageSize).Limit(form.PageSize).Find(&roleList)
	// 		data["list"] = roleList

	// 		global.App.DB.Model(&models.Role{}).Where("create_by = ?", managerId).Or("create_by IN (?)", managerSubIds).Count(&roleCount)
	// 		data["count"] = roleCount
	// 	} else {
	// 		global.App.DB.Where("create_by = ?", managerId).Offset((form.Page - 1) * form.PageSize).Limit(form.PageSize).Find(&roleList)
	// 		data["list"] = roleList

	// 		global.App.DB.Model(&models.Role{}).Where("create_by = ?", managerId).Count(&roleCount)
	// 		data["count"] = roleCount
	// 	}
	// }

	return
}

func (roleService *roleService) DoAdd(params request_admin.RoleAdd, c *gin.Context) (err error) {
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		value, _ := c.Keys["id"].(string)
		id, _ := strconv.ParseUint(value, 10, 64)

		role := models.Role{}
		role.Title = params.Title
		role.Description = params.Description
		role.Status = params.Status
		role.CreateBy = uint(id)

		err = global.App.DB.Create(&role).Error
		if err != nil {
			err = errors.New("admin.role.添加角色失败")
			return err
		}

		// 创建角色操作日志
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		err = CommonService.CreateAdminLog(uint(id), "角色模块", "创建角色", c.Request.Method, c.Request.URL.Path, utils.StructToJsonString(params), utils.StructToJsonString(role), ip, userAgent, 1, 0)
		if err != nil {
			return err
		}

		return nil
	})
}

// 编辑权限-获取权限
func (roleService *roleService) Edit(params request_admin.RoleEditDelete) (role *models.Role) {
	global.App.DB.Where("id = ?", params.ID).Find(&role)
	return
}

// 编辑权限-执行编辑
func (roleService *roleService) DoEdit(params request_admin.RoleEdit, c *gin.Context) (err error) {
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		value, _ := c.Keys["id"].(string)
		id, _ := strconv.ParseUint(value, 10, 64)

		//实例化permission
		role := models.Role{}

		global.App.DB.Where("id = ?", params.ID).Find(&role)
		role.Title = params.Title
		role.Description = params.Description
		role.Status = params.Status
		role.UpdateBy = uint(id)

		err := global.App.DB.Save(&role).Error
		if err != nil {
			err = errors.New("admin.role.编辑角色失败")
			return err
		}

		// 创建角色操作日志
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		err = CommonService.CreateAdminLog(uint(id), "角色模块", "编辑角色", c.Request.Method, c.Request.URL.Path, utils.StructToJsonString(params), utils.StructToJsonString(role), ip, userAgent, 1, 0)
		if err != nil {
			return err
		}

		return nil
	})
}

// 删除角色-执行删除
func (roleService *roleService) Delete(params request_admin.RoleEditDelete, c *gin.Context) (err error) {
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		if params.ID == 1 {
			err = errors.New(trans.Trans("admin.role.禁止删除超管角色"))
			return err
		}

		//删除角色
		err = global.App.DB.Delete(&models.Role{}, params.ID).Error
		if err != nil {
			err = errors.New(trans.Trans("admin.role.删除角色失败"))
			return err
		}

		//删除菜单授权
		err = global.App.DB.Where("role_id = ?", params.ID).Delete(&models.RolePermission{}).Error
		if err != nil {
			err = errors.New(trans.Trans("admin.permission.删除关联权限失败"))
			return err
		}

		//删除部门授权
		err = global.App.DB.Where("role_id = ?", params.ID).Delete(&models.RoleDept{}).Error
		if err != nil {
			err = errors.New(trans.Trans("admin.删除关联部门失败"))
			return err
		}

		// 创建权限操作日志
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		value, _ := c.Keys["id"].(string)
		id, _ := strconv.ParseUint(value, 10, 64)
		err = CommonService.CreateAdminLog(uint(id), "角色模块", "删除角色", c.Request.Method, c.Request.URL.Path, utils.StructToJsonString(params), "", ip, userAgent, 1, 0)
		if err != nil {
			return err
		}

		return nil
	})
}

// 角色菜单授权
func (roleService *roleService) PermissionAuth(params request_admin.RoleEditDelete) (role models.Role, permissionList []models.Permission) {
	id := params.ID

	global.App.DB.Find(&role, id)

	//获取所有权限列表
	global.App.DB.Where("module_id = ?", 0).Preload("PermissionItem.PermissionItem").Find(&permissionList)

	//获取当前角色拥有的权限,并把权限id放在一个map对象中
	rolePermission := []models.RolePermission{}
	global.App.DB.Where("role_id = ?", id).Find(&rolePermission)
	rolePermissionMap := make(map[uint]uint)
	for _, v := range rolePermission {
		rolePermissionMap[v.PermissionID] = v.PermissionID
	}

	//循环遍历所有权限数据,判断当前权限的id是否在角色权限的map对象中,如果是的话给当前数据加入checked属性
	for i := 0; i < len(permissionList); i++ { //循环权限列表
		if _, ok := rolePermissionMap[permissionList[i].ID.ID]; ok { // 判断当前权限是否在角色权限的map对象中
			permissionList[i].Checked = true
		}
		for j := 0; j < len(permissionList[i].PermissionItem); j++ { // 判断当前权限的子栏位是否在角色权限的map中
			if _, ok := rolePermissionMap[permissionList[i].PermissionItem[j].ID.ID]; ok { // 判断当前权限是否在角色权限的map对象中
				permissionList[i].PermissionItem[j].Checked = true
			}
			for jj := 0; jj < len(permissionList[i].PermissionItem[j].PermissionItem); jj++ { // 判断当前权限的子栏位是否在角色权限的map中
				if _, ok := rolePermissionMap[permissionList[i].PermissionItem[j].PermissionItem[jj].ID.ID]; ok { // 判断当前权限是否在角色权限的map对象中
					permissionList[i].PermissionItem[j].PermissionItem[jj].Checked = true
				}
			}
		}
	}

	return
}

// 角色菜单授权 - 执行授权
func (roleService *roleService) PermissionDoAuth(params request_admin.RolePermissionAuth, c *gin.Context) (err error) {
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		//获取表单提交的权限id切片
		roleID := params.RoleId
		permissionIds := params.PermissionNode

		// 如果切片中存在相同的元素，仅仅留下一个
		unique := make(map[uint]bool)
		for _, v := range permissionIds {
			unique[v] = true
		}
		permissionIds = []uint{}
		for k := range unique {
			permissionIds = append(permissionIds, k)
		}

		//先删除当前角色对应的权限
		rolePermission := []models.RolePermission{}
		global.App.DB.Where("role_id = ?", roleID).Delete(&rolePermission)

		//循环遍历permissionIds,增加当前角色对应的权限
		for _, v := range permissionIds {
			temp := models.RolePermission{}
			temp.RoleID = roleID
			temp.PermissionID = v
			rolePermission = append(rolePermission, temp)
		}
		global.App.DB.Create(&rolePermission)

		// 创建权限操作日志
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		value, _ := c.Keys["id"].(string)
		id, _ := strconv.ParseUint(value, 10, 64)

		err = CommonService.CreateAdminLog(uint(id), "角色模块", "角色授权", c.Request.Method, c.Request.URL.Path, utils.StructToJsonString(params), "", ip, userAgent, 1, 0)
		if err != nil {
			return err
		}

		return nil
	})
}

// 角色部门授权
func (roleService *roleService) DeptAuth(params request_admin.RoleEditDelete) (role models.Role, deptList []models.Dept) {
	id := params.ID

	global.App.DB.Find(&role, id)

	global.App.DB.Where("parent_id = ?", 0).Preload("DeptItem").Find(&deptList)

	roleDept := []models.RoleDept{}
	global.App.DB.Where("role_id = ?", id).Find(&roleDept)
	roleDeptMap := make(map[uint]uint)
	for _, v := range roleDept {
		roleDeptMap[v.DeptID] = v.DeptID
	}

	for i := 0; i < len(deptList); i++ {
		if _, ok := roleDeptMap[deptList[i].ID.ID]; ok {
			deptList[i].Checked = true
		}
		for j := 0; j < len(deptList[i].DeptItem); j++ {
			if _, ok := roleDeptMap[deptList[i].DeptItem[j].ID.ID]; ok {
				deptList[i].DeptItem[j].Checked = true
			}
		}
	}

	return
}

// 角色部门授权 - 执行授权
func (roleService *roleService) DeptDoAuth(params request_admin.RoleDeptAuth, c *gin.Context) (err error) {
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		roleID := params.RoleId
		deptIds := params.DeptNode

		// deptIds 去重
		unique := make(map[uint]bool)
		for _, v := range deptIds {
			unique[v] = true
		}
		deptIds = []uint{}
		parentIds := []uint{}
		for k := range unique {
			deptIds = append(deptIds, k)
			//判断当前部门是否存在上级部门，如果存在，则把上级部门加入到deptIds中
			dept := models.Dept{}
			global.App.DB.Where("id = ?", k).Select("id,parent_id").Find(&dept)
			if dept.ParentId != 0 && !utils.Contains(deptIds, dept.ParentId) {
				deptIds = append(deptIds, dept.ParentId)
				parentIds = append(parentIds, dept.ParentId)
			} else if dept.ParentId == 0 && !utils.Contains(parentIds, dept.ID.ID) {
				parentIds = append(parentIds, dept.ID.ID)
			}
		}

		// fmt.Printf("deptIds:%+v\n", deptIds)
		// fmt.Printf("parentIds:%+v\n", parentIds)

		roleDept := []models.RoleDept{}
		// 查询所有父级部门，根据父级部门的id查询子级部门id,然后删除所有父级和子级部门的角色授权
		subIds := []uint{}
		if len(parentIds) > 0 {
			global.App.DB.Model(&models.Dept{}).Where("parent_id IN (?)", parentIds).Select("id").Find(&subIds)
		}

		ids := append(parentIds, subIds...)
		global.App.DB.Where("role_id = ?", roleID).Where("dept_id IN (?)", ids).Delete(&roleDept)
		for _, v := range deptIds {
			temp := models.RoleDept{}
			temp.RoleID = roleID
			temp.DeptID = v
			roleDept = append(roleDept, temp)
		}
		global.App.DB.Create(&roleDept)

		// 创建部门操作日志
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		value, _ := c.Keys["id"].(string)
		id, _ := strconv.ParseUint(value, 10, 64)

		err = CommonService.CreateAdminLog(uint(id), "角色模块", "角色授权", c.Request.Method, c.Request.URL.Path, utils.StructToJsonString(params), "", ip, userAgent, 1, 0)
		if err != nil {
			return err
		}

		return nil
	})
}
