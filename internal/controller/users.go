package controller

import (
	"net/http"
	"videohub/internal/service"

	"github.com/gin-gonic/gin"
)

type Users struct {
	userAvatarService *service.User_avatar
	userListService   *service.User_list
	userService       *service.User
}

func NewUsers(uas *service.User_avatar, uls *service.User_list, us *service.User) *Users {
	return &(Users{userAvatarService: uas, userListService: uls, userService: us})
}

/*
*@author:廖嘉鹏
*@create_at:2024/10/17
 */
// 测试，这是一个样版
func (uc *Users) User_test(c *gin.Context) {
	//调用service里的方法进行处理
	uc.userAvatarService.Test()
	uc.userListService.Test()
	uc.userService.Test()
	c.JSON(http.StatusOK, gin.H{"message": "user"})
}

//一个API对应一个函数
//追加在下面
