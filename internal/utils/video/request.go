package video

import "mime/multipart"

type UploadChunkRequest struct {
	UploadID  string                `form:"upload_id" binding:"required"`
	ChunkData *multipart.FileHeader `form:"chunk_data" binding:"required"`
	ChunkID   int                   `form:"chunk_id" binding:"required"`
	ChunkSize int                   `form:"chunk_size" binding:"required"`
	ChunkHash string                `form:"chunk_hash" binding:"required"`
}

type CompleteUploadRequest struct {
	UploadID    string                `form:"upload_id" binding:"required"`
	ChunkEndID  int                   `form:"chunk_end_id" binding:"required"`
	Title       string                `form:"title" binding:"required"`
	Description string                `form:"description" binding:"required"`
	Cover       *multipart.FileHeader `form:"cover" binding:"required"`
	VideoHash   string                `form:"video_hash" binding:"required"`
	UploaderID  uint                  `form:"uploader_id" binding:"required"`
}
