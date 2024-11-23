package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	ID   uint64 `json:"id"`
	Role uint8  `json:"role"`
}

type MyCustomClaims struct {
	Payload   Payload
	IssuedAt  int64 `json:"iat"`
	ExpiresAt int64 `json:"exp"`
	jwt.RegisteredClaims
}

func GenerateJWT(pyaload Payload, key string, duration uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		Payload:   pyaload,
		IssuedAt:  time.Now().UnixMilli(),
		ExpiresAt: time.Now().Add(time.Second * time.Duration(duration)).UnixMilli(),
	})
	ss, err := token.SignedString([]byte(key))
	return "Bearer " + ss, err
}

func ParseJWT(tokenString string, key string) (Payload, error) {
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return Payload{}, err
	} else if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		if time.Now().UnixMilli() > claims.ExpiresAt {
			return Payload{}, jwt.ErrTokenExpired
		}
		return claims.Payload, nil
	} else {
		return Payload{}, jwt.ErrTokenInvalidClaims
	}
}

// SendEmailVerification 发送邮箱验证
func SendEmailVerification(email string) error {
	// 模拟发送验证码到邮箱
	fmt.Printf("发送验证码到邮箱: %s\n", email)
	return nil
}

// generateSalt 生成随机盐值
func GenerateSalt(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

// hashPassword 对密码和盐值进行哈希
func HashPassword(password string, salt string) string {
	saltedPassword := password + salt
	hash := sha256.Sum256([]byte(saltedPassword))
	return hex.EncodeToString(hash[:])
}

func GetUserID(c *gin.Context) (uint64, error) {
	idValue, exists := c.Get("id")
	if !exists {
		return 0, errors.New("上下文中不存在用户 ID")
	}
	id, ok := idValue.(uint64)
	if !ok {
		return 0, errors.New("用户 ID 类型错误")
	}
	return id, nil
}

func CheckFile(file *multipart.FileHeader, types []string, maxSize int64) error {
	if file.Size > maxSize {
		return errors.New("文件大小超过限制")
	}

	fileExt := filepath.Ext(file.Filename)
	for _, t := range types {
		if t == fileExt {
			return nil
		}
	}

	return errors.New("文件类型不支持")
}

func SaveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return errors.New("文件打开失败")
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return errors.New("创建文件夹失败")
	}

	out, err := os.Create(dst)
	if err != nil {
		return errors.New("创建文件失败")
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return errors.New("拷贝文件失败")
	}
	return nil
}
