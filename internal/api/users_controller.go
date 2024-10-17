package api

import (
	"net/http"
	"videohub/internal/service"

	"github.com/gin-gonic/gin"
)

type Users_controller struct {
	user_avatar_service *service.User_avatar_service
	user_list_service   *service.User_list_service
	user_service        *service.User_service
}

func NewUsers_controller(uas *service.User_avatar_service, uls *service.User_list_service, us *service.User_service) *Users_controller {
	return &(Users_controller{user_avatar_service: uas, user_list_service: uls, user_service: us})
}

/*
*@author:廖嘉鹏
*@create_at:2024/10/17
 */
// 测试，这是一个样版
func (uc *Users_controller) User_test(c *gin.Context) {
	//调用service里的方法进行处理
	uc.user_avatar_service.Test()
	uc.user_list_service.Test()
	uc.user_service.Test()
	c.JSON(http.StatusOK, gin.H{"message": "用户首页"})
}

//一个API对应一个函数
//追加在下面
