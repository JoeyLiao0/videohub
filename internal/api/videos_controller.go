package api

import (
	"net/http"
	"videohub/internal/service"

	"github.com/gin-gonic/gin"
)

type Videos_controller struct {
	video_service        *service.Video_service
	video_list_service   *service.Video_list_service
	video_upload_service *service.Video_upload_service
}

func NewVideos_controller(vs *service.Video_service, vls *service.Video_list_service, vus *service.Video_upload_service) *Videos_controller {
	return &(Videos_controller{video_service: vs, video_list_service: vls, video_upload_service: vus})
}

/*
*@author:廖嘉鹏
*@create_at:2024/10/17
 */
// 测试，这是一个样版
func (vc *Videos_controller) Video_test(c *gin.Context) {
	//调用service里的方法进行处理
	vc.video_upload_service.Test()
	vc.video_list_service.Test()
	vc.video_service.Test()
	c.JSON(http.StatusOK, gin.H{"message": "用户首页"})
}

//一个API对应一个函数
//追加在下面
