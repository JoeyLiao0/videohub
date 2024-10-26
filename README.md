
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
│  ├─ config
│  │  └─ config.go
│  ├─ controller
│  │  ├─ users.go
│  │  └─ videos.go
│  ├─ model
│  │  ├─ collection.go
│  │  ├─ comment.go
│  │  ├─ user.go
│  │  └─ video.go
│  ├─ repository
│  │  ├─ collection.go
│  │  ├─ comment.go
│  │  ├─ user.go
│  │  └─ video.go
│  ├─ router
│  │  └─ router.go
│  ├─ service
│  │  ├─ user.go
│  │  ├─ user_avatar.go
│  │  ├─ user_list.go
│  │  ├─ video.go
│  │  ├─ video_list.go
│  │  └─ video_upload.go
│  └─ utils
└─ storage
   ├─ images
   └─ videos
      ├─ cover
      ├─ data
      └─ tmp

```