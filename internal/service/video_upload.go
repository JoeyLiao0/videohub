package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"videohub/config"
	"videohub/internal/model"
	"videohub/internal/repository"
)

type VideoUpload struct {
	videoRepo *repository.Video
}

func NewVideoUpload(vr *repository.Video) *VideoUpload {
	return &VideoUpload{videoRepo: vr}
}

/*
*@author:李逸城
*@create_at:2024/11/7
 */

func (vus *VideoUpload) HandleVideoChunk(videoChunk model.VideoChunk, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("failed to open file: %v", err)
		return err
	}
	// defer最后关闭文件
	defer file.Close()
	fileSize := fileHeader.Size

	// 验证切片大小(byte)
	if fileSize != int64(videoChunk.ChunkSize) {
		err := fmt.Errorf("chunk size mismatch: expected %d, got %d", videoChunk.ChunkSize, fileSize)
		log.Println(err)
		return err
	}

	// 验证哈希值(sha256)
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		log.Printf("failed to calculate hash: %v", err)
		return err
	}
	calculatedHash := hex.EncodeToString(hasher.Sum(nil))
	if calculatedHash != videoChunk.ChunkHash {
		err := fmt.Errorf("chunk hash mismatch: expected %s, got %s", videoChunk.ChunkHash, calculatedHash)
		log.Println(err)
		return err
	}

	// 从配置中获取临时路径
	tmpDir := config.AppConfig.Storage.VideosChunk
	saveDir := filepath.Join(tmpDir, videoChunk.UploadID)
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		log.Printf("failed to create directory: %v", err)
		return err
	}

	// 创建切片文件路径(/tmp/{uploadID}/{uploadID}_{chunkID}.xxx)
	tempSavePath := filepath.Join(saveDir, fmt.Sprintf("%s_%d_tmp", videoChunk.UploadID, videoChunk.ChunkID))
	log.Printf("save tempSaveFile at: %s", tempSavePath)
	out, err := os.Create(tempSavePath)
	if err != nil {
		log.Printf("failed to save chunk: %v", err)
		return err
	}
	defer out.Close()

	// 重置文件指针
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		log.Printf("failed to seek file: %v", err)
		return err
	}

	// 调用DAO层写入视频切片
	if err := vus.videoRepo.Save(file, tempSavePath); err != nil {
		return fmt.Errorf("failed to save chunk: %v", err)
	}

	log.Printf("Chunk saved successfully: %s", tempSavePath) // 日志
	return nil
}

// HandleVideoComplete 处理组合完整视频逻辑
func (vus *VideoUpload) HandleVideoComplete(video model.Video, chunkEndID int, coverFile multipart.File, videoHash string) error {
	uploadIDStr := strconv.FormatInt(video.UploadID, 10)

	// 调用DAO层获取视频切片列表（[]string）
	chunks, err := vus.videoRepo.GetVideoChunksByUploadID(uploadIDStr, chunkEndID)
	if err != nil {
		log.Printf("Error retrieving video chunks for upload ID %s: %v", uploadIDStr, err)
		return err
	}

	// 合并切片文件到输出视频文件
	videoPath := filepath.Join(config.AppConfig.Storage.VideosData, fmt.Sprintf("%d.mp4", video.UploadID))
	if err := vus.videoRepo.SaveChunks(chunks, videoPath); err != nil {
		log.Printf("Error saving chunks to output video file %s: %v", videoPath, err)
		return err
	}

	// 保存完整视频路径
	video.VideoPath = videoPath

	// 计算合并后视频的 SHA-256 哈希
	file, err := os.Open(videoPath)
	if err != nil {
		log.Printf("Error opening  video file %s: %v", videoPath, err)
		return err
	}
	// defer file.Close()
	// 使用defer后面删除切片会无法关闭

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Printf("Error copying  video file %s: %v", videoPath, err)
		return err
	}

	file.Close()

	hashValue := hex.EncodeToString(hash.Sum(nil))

	// 校验哈希
	if hashValue != videoHash {
		file.Close()
		// 若哈希校验错误，删除合并后的视频文件
		if err := os.Remove(videoPath); err != nil {
			log.Printf("Error deleting video file %s: %v", videoPath, err)
			return err
		}
		log.Printf("hash mismatch: expected %s, got %s", videoHash, hashValue)
		log.Printf("Video file %s successfully deleted after hash verification", videoPath)
		return err
	}
	log.Printf("Hash verification succeeded for video %s", videoPath)

	// 创建封面文件路径 (/cover/{uploadID}.png)
	coverPath := filepath.Join(config.AppConfig.Storage.VideosCover, fmt.Sprintf("%s.png", uploadIDStr))
	coverOutFile, err := os.Create(coverPath)
	if err != nil {
		log.Printf("Error creating cover file at %s: %v", coverPath, err)
		return err
	}
	defer coverOutFile.Close()

	// 调用DAO层写入视频封面文件
	if err := vus.videoRepo.Save(coverFile, coverPath); err != nil {
		return fmt.Errorf("failed to save cover file: %v", err)
	}

	// 更新视频结构体中的封面路径
	video.CoverPath = coverPath

	// 保存完整视频路径和封面路径
	if err := vus.videoRepo.SaveCompleteVideo(video); err != nil {
		log.Printf("Error saving complete video metadata for upload ID %s: %v", uploadIDStr, err)
		return err
	}

	// 调用DAO层删除所有切片文件
	vus.videoRepo.DeleteChunks(int(video.UploadID))
	log.Printf("Chunk files for upload ID %s successfully deleted", uploadIDStr)

	log.Printf("Video upload %s successfully completed and saved at %s with cover at %s", uploadIDStr, videoPath, coverPath)
	return nil
}
