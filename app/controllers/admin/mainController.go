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
	"shalabing-gin/app/common/response"
	"shalabing-gin/app/services/services_admin"

	"github.com/gin-gonic/gin"
)

type MainController struct{}

func (con MainController) Index(c *gin.Context) {
	data, err := services_admin.MainService.Index(c)
	if err != nil {
		response.TokenFail(c)
		return
	}
	response.Success(c, data)
}
