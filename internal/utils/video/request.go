package video

import "mime/multipart"

type GetVideosRequest struct {
	Status *int   `form:"status"`
	Like   string `form:"like"`
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

type UpdateVideoStatusRequest struct {
	NewStatus int8 `json:"new_status" binding:"required"`
}

type AddCommentRequest struct {
	UserID          uint   `json:"user_id" binding:"required"`
	CommentContent  string `json:"comment" binding:"required"`
	FatherCommentID int   `json:"father_comment_id" binding:"required"`
}

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
