package middleware

import (
	"encoding/json"
	"shalabing-gin/app/common/response"
	"shalabing-gin/app/models"
	"shalabing-gin/app/services"
	"shalabing-gin/global"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTAuth(GuardName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")
		if tokenStr == "" {
			response.TokenFail(c)
			c.Abort()
			return
		}
		tokenStr = tokenStr[len(services.TokenType)+1:]

		// Token 解析校验
		token, err := jwt.ParseWithClaims(tokenStr, &services.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(global.App.Config.Jwt.Secret), nil
		})
		if err != nil || services.JwtService.IsInBlacklist(tokenStr) {
			response.TokenFail(c)
			c.Abort()
			return
		}

		claims := token.Claims.(*services.CustomClaims)
		// Token 发布者校验
		if claims.Issuer != GuardName {
			response.TokenFail(c)
			c.Abort()
			return
		}

		// token 续签
		// claims.ExpiresAt-time.Now().Unix()  令牌剩余时长
		// global.App.Config.Jwt.RefreshGracePeriod  令牌宽限时长
		// fmt.Println(claims.ExpiresAt-time.Now().Unix(), global.App.Config.Jwt.RefreshGracePeriod)
		if claims.ExpiresAt-time.Now().Unix() < global.App.Config.Jwt.RefreshGracePeriod {
			// 在此期间内的token可以正常使用，锁过期重新颁发token
			lock := global.Lock("refresh_token_lock", global.App.Config.Jwt.JwtBlacklistGracePeriod)
			if lock.Get() {
				user, err := services.JwtService.GetUserInfo(GuardName, claims.Id)
				if err != nil {
					global.App.Log.Error(err.Error())
					lock.Release()
				} else {
					tokenData, _, _ := services.JwtService.CreateToken(GuardName, user)
					c.Header("new-token", tokenData.AccessToken)
					c.Header("new-expires-in", strconv.Itoa(tokenData.ExpiresIn))
					_ = services.JwtService.JoinBlackList(token)
				}
			}
		}

		c.Set("token", token)
		c.Set("id", claims.Id)
		switch GuardName {
		case services.AdminGuardName:
			manager := models.Manager{}
			_ = global.App.DB.Where("id = ?", claims.Id).First(&manager).Error
			c.Set("username", manager.Username)
			c.Set("nickname", manager.Nickname)
			c.Set("is_super", manager.IsSuper)
			c.Set("dept_id", manager.DeptId)
			userinfoSlice, _ := json.Marshal(manager)
			c.Set("userinfo", string(userinfoSlice))
		case services.ApiGuardName:
			user := models.User{}
			_ = global.App.DB.Where("id = ?", claims.Id).First(&user).Error
			c.Set("nickname", user.Username)
		}
	}
}
