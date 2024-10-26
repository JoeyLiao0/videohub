package controller

import (
	"net/http"
	"videohub/internal/service"

	"github.com/gin-gonic/gin"
)

type Videos struct {
	videoService       *service.Video
	videoListService   *service.Video_list
	videoUploadService *service.Video_upload
}

func NewVideos(vs *service.Video, vls *service.Video_list, vus *service.Video_upload) *Videos {
	return &(Videos{videoService: vs, videoListService: vls, videoUploadService: vus})
}

/*
*@author:廖嘉鹏
*@create_at:2024/10/17
 */
// 测试，这是一个样版
func (vc *Videos) Video_test(c *gin.Context) {
	//调用service里的方法进行处理
	vc.videoService.Test()
	vc.videoListService.Test()
	vc.videoUploadService.Test()
	c.JSON(http.StatusOK, gin.H{"message": "video"})
}

//一个API对应一个函数
//追加在下面
