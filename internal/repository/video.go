package repository

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"videohub/config"
	"videohub/internal/model"

	"gorm.io/gorm"
)

type Video struct {
	DB *gorm.DB
}

func NewVideo(db *gorm.DB) *Video {
	return &Video{DB: db}
}

// DeleteChunks 将多个切片文件删除
func (v *Video) DeleteChunks(UploadID int) error {
	saveDir := fmt.Sprintf("/tmp/%d", UploadID)
	if err := os.RemoveAll(saveDir); err != nil {
		log.Printf("Error deleting chunk files for upload ID %d: %v", UploadID, err)
		return err
	}
	return nil
}

// SaveChunks 将多个切片合并到一个输出文件
func (v *Video) SaveChunks(chunkPaths []string, outputPath string) error {
	// 创建输出文件
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %v", outputPath, err)
	}
	defer outFile.Close()

	// 逐个追加每个切片文件内容
	for _, chunkPath := range chunkPaths {
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			return fmt.Errorf("failed to open chunk file %s: %v", chunkPath, err)
		}

		// 将切片文件内容复制到输出文件
		if _, err := io.Copy(outFile, chunkFile); err != nil {
			chunkFile.Close()
			return fmt.Errorf("failed to copy chunk file %s to output file %s: %v", chunkPath, outputPath, err)
		}
		chunkFile.Close()
	}
	return nil
}

// Save 根据指定的类型保存视频切片或封面文件
func (v *Video) Save(file io.Reader, savePath string) error {
	// 创建保存文件
	out, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create file at %s: %v", savePath, err)
	}
	defer out.Close()

	// 写入内容
	if _, err := io.Copy(out, file); err != nil {
		os.Remove(savePath) // 如果写入失败，删除不完整的文件
		return fmt.Errorf("failed to save file at %s: %v", savePath, err)
	}

	return nil
}

// GetVideoChunksByUploadID 从 /tmp/{uploadID} 目录中找到所有切片文件，按编号排序返回切片路径列表
func (repo *Video) GetVideoChunksByUploadID(uploadID string, chunkEndID int) ([]string, error) {
	// 获取配置中的切片目录路径
	dirPath := filepath.Join(config.AppConfig.Storage.VideosChunk, uploadID)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var chunks []string
	for _, entry := range entries {
		chunkPath := filepath.Join(dirPath, entry.Name())
		chunks = append(chunks, chunkPath)
	}

	// 把各个切片文件排序
	sort.Slice(chunks, func(i, j int) bool {
		id1, _ := strconv.Atoi(filepath.Base(chunks[i])[len(uploadID)+1:])
		id2, _ := strconv.Atoi(filepath.Base(chunks[j])[len(uploadID)+1:])
		return id1 < id2
	})

	// 检查是否有缺失切片
	if len(chunks) < chunkEndID {
		return nil, errors.New("chunks are missing")
	}
	return chunks, nil
}

// SaveCompleteVideo 保存完整视频到数据库
func (vr *Video) SaveCompleteVideo(video model.Video) error {
	// GORM创建新纪录
	return vr.DB.Create(&video).Error
}

// UpdateVideoStatus 更新视频状态
func (vr *Video) UpdateVideoStatus(id int64, newStatus int8) error {
	result := vr.DB.Model(&model.Video{}).Where("upload_id = ?", id).Update("video_status", newStatus)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("视频未找到")
	}
	return nil
}

// 查询视频列表
func (vr *Video) FindVideos(status *int, like *string, offset, limit int) ([]model.Video, int64, error) {
	var videos []model.Video
	var total int64

	query := vr.DB.Model(&model.Video{})

	// 状态筛选
	if status != nil {
		query = query.Where("video_status = ?", *status)
	}

	// 标题模糊搜索
	if like != nil {
		query = query.Where("title LIKE ?", "%"+*like+"%")
	}

	// 分页查询
	err := query.Count(&total).Offset(offset).Limit(limit).Find(&videos).Error
	if err != nil {
		return nil, 0, err
	}

	return videos, total, nil
}