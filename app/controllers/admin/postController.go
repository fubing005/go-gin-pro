/*
*

		_ooOoo_
	              o8888888o
	              88" . "88
	              (| -_- |)
	              O\  =  /O
	           ____/`---'\____
	         .'  \\|     |//  `.
	        /  \\|||  :  |||//  \
	       /  _||||| -:- |||||_  \
	       |   | \\\  -  /'| |   |
	       | \_|  `\`---'//  |_/ |
	       \  .-\__ `-. -'__/-.  /
	     ___`. .'  /--.--\  `. .'___
	  ."" '<  `.___\_<|>_/___.' _> \"".
	 | | :  `- \`. ;`. _/; .'/ /  .' ; |
	 \  \ `-.   \_\_`. _.'_/_/  -' _.' /

===================================

	佛祖加持 代码无BUG

===================================
*/
package admin

import (
	"shalabing-gin/app/common/request"
	"shalabing-gin/app/common/request/request_admin"
	"shalabing-gin/app/common/response"
	"shalabing-gin/app/services/services_admin"

	"github.com/gin-gonic/gin"
)

type PostController struct{}

// 通用数据
func (con PostController) Common(c *gin.Context) {
	commonData := services_admin.PostService.Common()
	response.Success(c, commonData)
}

// 获取所有岗位
func (con PostController) Index(c *gin.Context) {
	var form request.PageQuery
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	postList := services_admin.PostService.GetPosts(form)
	response.Success(c, postList)
}

// 添加岗位-执行添加
func (con PostController) DoAdd(c *gin.Context) {
	var form request_admin.PostAdd
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	err := services_admin.PostService.DoAdd(form, c)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// 编辑岗位-获取单个岗位、顶级岗位
func (con PostController) Edit(c *gin.Context) {
	var form request_admin.PostEditDelete
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	post := services_admin.PostService.Edit(form)
	response.Success(c, post)
}

// 编辑岗位-执行编辑
func (con PostController) DoEdit(c *gin.Context) {
	var form request_admin.PostEdit
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	err := services_admin.PostService.DoEdit(form, c)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// 删除岗位
func (con PostController) Delete(c *gin.Context) {
	var form request_admin.PostEditDelete
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	err := services_admin.PostService.Delete(form, c)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, nil)
}
