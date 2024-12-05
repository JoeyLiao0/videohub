package service

import (
	"fmt"
	"net/http"
	"path/filepath"
	"videohub/config"
	"videohub/internal/model"
	"videohub/internal/repository"
	"videohub/internal/utils"
	"videohub/internal/utils/video"

	"github.com/sirupsen/logrus"
)

type VideoUpload struct {
	videoRepo *repository.Video
}

func NewVideoUpload(vr *repository.Video) *VideoUpload {
	return &VideoUpload{videoRepo: vr}
}

func (vus *VideoUpload) HandleVideoChunk(request *video.UploadChunkRequest) *utils.Response {
	if err := utils.CheckFile(request.ChunkData, []string{".mp4", ".avi", ".mov", ".mkv"}, 32<<20); err != nil {
		logrus.Debug(err.Error())
		return utils.Error(http.StatusBadRequest, "文件格式错误或文件过大")
	}
	fileSize := request.ChunkData.Size

	// 验证切片大小(byte)
	if fileSize != int64(request.ChunkSize) {
		logrus.Debugf("chunk size mismatch: expected %d, got %d", request.ChunkSize, fileSize)
		return utils.Error(http.StatusBadRequest, "分片大小不匹配")
	}

	hashValue, err := utils.CalculateFileHash(request.ChunkData)
	if err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	if hashValue != request.ChunkHash {
		logrus.Debugf("chunk hash mismatch: expected %s, got %s", request.ChunkHash, hashValue)
		return utils.Error(http.StatusBadRequest, "分片哈希不匹配")
	}

	// 创建切片文件路径(/tmp/{uploadID}/{uploadID}_{chunkID}.xxx)
	tmpDir := config.AppConfig.Storage.VideosChunk
	saveDir := filepath.Join(tmpDir, string(request.UploadID))
	tempSavePath := filepath.Join(saveDir, fmt.Sprintf("%s_%d%s", request.UploadID, request.ChunkID, filepath.Ext(request.ChunkData.Filename)))

	if err := utils.SaveFile(request.ChunkData, tempSavePath); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	logrus.Debug("Video chunk Upload successfully")
	return utils.Success(http.StatusOK)
}

// HandleVideoComplete 处理组合完整视频逻辑
func (vus *VideoUpload) HandleVideoComplete(request *video.CompleteUploadRequest) *utils.Response {
	// 检查文件的类型、大小
	if err := utils.CheckFile(request.Cover, []string{".png", ".jpg", ".jpeg"}, 8<<20); err != nil {
		logrus.Debug(err.Error())
		return utils.Error(http.StatusBadRequest, "文件格式错误或文件过大")
	}

	// 调用DAO层获取视频切片列表（[]string）
	dirPath := filepath.Join(config.AppConfig.Storage.VideosChunk, request.UploadID)
	chunks, err := utils.ListFilesSortedByName(dirPath, request.ChunkEndID)
	if err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	// 计算合并后视频的 SHA-256 哈希
	hashValue, err := utils.CalculateFileHash(chunks)
	if err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	// 校验哈希
	if hashValue != request.VideoHash {
		logrus.Debugf("hash mismatch: expected %s, got %s", request.VideoHash, hashValue)
		return utils.Error(http.StatusBadRequest, "哈希校验错误")
	}

	// 创建封面文件路径 (/cover/{uploadID}.ext)
	coverExt := filepath.Ext(request.Cover.Filename)
	coverPath := filepath.Join(config.AppConfig.Storage.VideosCover, fmt.Sprintf("%s%s", request.UploadID, coverExt))
	if err := utils.SaveFile(request.Cover, coverPath); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	// 合并切片文件到输出视频文件
	videoPath := filepath.Join(config.AppConfig.Storage.VideosData, fmt.Sprintf("%s%s", request.UploadID, filepath.Ext(chunks[0])))
	if err := utils.MergeFiles(chunks, videoPath); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	newVideo := model.Video{
		UploadID:     request.UploadID,
		Title:        request.Title,
		Description:  request.Description,
		CoverPath:    utils.GetURLPath(config.AppConfig.Static.Cover, fmt.Sprintf("%s%s", request.UploadID, coverExt)),
		VideoPath:    utils.GetURLPath(config.AppConfig.Static.Video, fmt.Sprintf("%s%s", request.UploadID, filepath.Ext(chunks[0]))),
		UploaderName: request.UploaderName,
	}

	// 保存完整视频路径和封面路径
	if err := vus.videoRepo.CreateVideo(&newVideo); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	// 调用DAO层删除所有切片文件
	chunkDir := filepath.Join(config.AppConfig.Storage.VideosChunk, request.UploadID)
	if err := utils.RemoveDir(chunkDir); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	logrus.Debugf("Video upload %s successfully completed and saved at %s with cover at %s", request.UploadID, videoPath, coverPath)
	return utils.Success(http.StatusOK)
}
