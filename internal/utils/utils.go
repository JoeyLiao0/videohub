package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"text/template"
	"time"
	"videohub/config"
	"videohub/global"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

// Response 响应结构体
type Response struct {
	StatusCode int         `json:"code"`
	ErrorMsg   string      `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

// Success 返回成功响应, 不带数据
func Success(statusCode int) *Response {
	return &Response{StatusCode: statusCode}
}

// Ok 返回成功响应, 带数据
func Ok(statusCode int, data interface{}) *Response {
	return &Response{StatusCode: statusCode, Data: data}
}

// Error 返回错误响应
func Error(statusCode int, errorMsg string) *Response {
	return &Response{StatusCode: statusCode, ErrorMsg: errorMsg}
}

// Payload JWT载荷
type Payload struct {
	ID   uint `json:"id"`
	Role int8 `json:"role"`
}

// MyCustomClaims 自定义声明
type MyCustomClaims struct {
	Payload   Payload
	IssuedAt  int64 `json:"iat"`
	ExpiresAt int64 `json:"exp"`
	jwt.RegisteredClaims
}

// GenerateJWT 生成JWT
func GenerateJWT(pyaload Payload, key string, duration uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		Payload:   pyaload,
		IssuedAt:  time.Now().UnixMilli(),
		ExpiresAt: time.Now().Add(time.Second * time.Duration(duration)).UnixMilli(),
	})
	ss, err := token.SignedString([]byte(key))
	return "Bearer " + ss, err
}

// ParseJWT 解析JWT
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

// GenerateCode 生成随机验证码
func GenerateCode(length int) string {
	const digits = "0123456789"
	code := make([]byte, length)
	for i := range code {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			log.Fatal(err.Error())
		}
		code[i] = digits[randomIndex.Int64()]
	}
	log.Println(string(code))
	return string(code)
}

// LoadAndFillTemplate 加载并填充 email 模板
func LoadAndFillTemplate(filePath string, data interface{}) (string, error) {
	templateFile, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("email").Parse(string(templateFile))
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// SendEmailVerification 发送邮箱验证
func SendEmailVerification(to string, code string) error {
	username := config.AppConfig.Email.Username
	password := config.AppConfig.Email.Password
	host := config.AppConfig.Email.Host
	port := config.AppConfig.Email.Port
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("VideoHub <%s>", username))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Verification Code")
	data, err := LoadAndFillTemplate("template/email.html",
		map[string]interface{}{"verification_code": code, "expiration_time": config.AppConfig.Email.Expiration})
	if err != nil {
		log.Println(err.Error())
		return errors.New("加载模板文件失败")
	}
	m.SetBody("text/html", data)

	d := gomail.NewDialer(host, port, username, password)
	if err := d.DialAndSend(m); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// VerifyEmailCode 验证邮箱验证码
func VerifyEmailCode(email string, code string) error {
	c, err := global.Rdb.Get(global.Ctx, email).Result()
	if err != nil {
		return err
	}
	if c != code {
		return errors.New("code is invalid")
	}

	global.Rdb.Del(global.Ctx, email)
	return nil
}

// generateSalt 生成随机盐值
func GenerateSalt(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err.Error())
	}
	return hex.EncodeToString(bytes)
}

// hashPassword 对密码和盐值进行哈希
func HashPassword(password string, salt string) string {
	saltedPassword := password + salt
	hash := sha256.Sum256([]byte(saltedPassword))
	return hex.EncodeToString(hash[:])
}

// CheckFile 检查文件类型和大小
func CheckFile(file *multipart.FileHeader, types []string, maxSize int64) error {
	if file.Size > maxSize {
		return errors.New("file size exceeds the limit")
	}

	fileExt := filepath.Ext(file.Filename)
	for _, t := range types {
		if t == fileExt {
			return nil
		}
	}

	return errors.New("file type is not supported")
}

// SaveFile 保存文件
func SaveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return err
	}
	return nil
}

// MergeFiles 合并多个分片文件
func MergeFiles(filePaths []string, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	for _, filePath := range filePaths {
		src, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer src.Close()

		_, err = io.Copy(out, src)
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveDir 删除指定目录下的所有文件和子目录
func RemoveDir(dirPath string) error {
	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return err
	}
	// 删除目录下的所有文件和子目录
	if err := os.RemoveAll(dirPath); err != nil {
		return err
	}
	return nil
}

// RemoveFile 删除指定文件
func RemoveFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return err
	}
	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}

// CalculateFileHash 计算文件哈希值
func CalculateFileHash(input interface{}) (string, error) {
	hasher := sha256.New()

	switch v := input.(type) {
	case *multipart.FileHeader:
		// 处理 *multipart.FileHeader 类型
		file, err := v.Open()
		if err != nil {
			return "", err
		}
		defer file.Close()

		if _, err := io.Copy(hasher, file); err != nil {
			return "", err
		}

	case []string:
		// 处理 []string 类型
		for _, filePath := range v {
			logrus.Debug(filePath)
			file, err := os.Open(filePath)
			if err != nil {
				return "", err
			}
			defer file.Close()

			if _, err := io.Copy(hasher, file); err != nil {
				return "", err
			}
		}

	default:
		return "", fmt.Errorf("unsupported input type")
	}

	hash := hasher.Sum(nil)
	hashString := hex.EncodeToString(hash)

	return hashString, nil
}

// ListFilesSortedByName 列出指定目录下按名称排序的文件名
func ListFilesSortedByName(dirPath string, count int, prefix string) ([]string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}
	logrus.Debug(fileNames)
	sort.Slice(fileNames, func(i, j int) bool {
		id1, _ := strconv.Atoi(filepath.Base(fileNames[i])[len(prefix)+1:])
		id2, _ := strconv.Atoi(filepath.Base(fileNames[j])[len(prefix)+1:])
		return id1 < id2
	})

	for i, fileName := range fileNames {
		fileNames[i] = filepath.Join(dirPath, fileName)
	}

	// 检查是否有缺失
	if len(fileNames) < count {
		return nil, errors.New("files are missing")
	}
	return fileNames, nil
}

// GenerateUsername 生成随机用户名
func GenerateUsername(n int) (string, error) {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		result[i] = letters[num.Int64()]
	}
	return string(result), nil
}

// GetURLPath 获取静态资源路径
func GetURLPath(dir, path string) string {
	return fmt.Sprintf("%s%s/%s", config.AppConfig.Static.Base, dir, path)
}
