
```
videohub
├─ cmd
│  └─ main.go
├─ config
│  └─ application.yaml
├─ Dockerfile
├─ go.mod
├─ go.sum
├─ internal
│  ├─ api
│  │  ├─ users_controller.go
│  │  └─ videos_controller.go
│  ├─ model
│  │  ├─ collection_model.go
│  │  ├─ comment_model.go
│  │  ├─ user_model.go
│  │  └─ video_model.go
│  ├─ repository
│  │  ├─ collection_repository.go
│  │  ├─ comment_repository.go
│  │  ├─ user_repository.go
│  │  └─ video_repository.go
│  ├─ service
│  │  ├─ user_avatar_service.go
│  │  ├─ user_list_service.go
│  │  ├─ user_service.go
│  │  ├─ video_list_service.go
│  │  ├─ video_service.go
│  │  └─ video_upload_service.go
│  └─ utils
│     └─ configTool.go
├─ README.md
└─ storage
   ├─ images
   └─ videos
      ├─ cover
      ├─ data
      └─ tmp

```