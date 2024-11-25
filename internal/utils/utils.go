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
	"text/template"
	"time"
	"videohub/config"
	"videohub/global"

	"github.com/golang-jwt/jwt/v5"
	"gopkg.in/gomail.v2"
)

type Response struct {
	StatusCode int         `json:"-"`
	ErrorMsg   string      `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Status     bool        `json:"status,omitempty"`
}

func Success(statusCode int) *Response {
	return &Response{StatusCode: statusCode, Status: true}
}

func Fail(statusCode int) *Response {
	return &Response{StatusCode: statusCode, Status: false}
}

func Ok(statusCode int, data interface{}) *Response {
	return &Response{StatusCode: statusCode, Data: data, Status: true}
}

func Error(statusCode int, errorMsg string) *Response {
	return &Response{StatusCode: statusCode, ErrorMsg: errorMsg, Status: false}
}


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

func GenerateCode(length int) string {
	// randSource := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	// code := fmt.Sprintf("%06d", randSource.Intn(1000000))
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
	// addr := fmt.Sprintf("%s:%d", host, port)
	// e := email.NewEmail()
	// e.From = fmt.Sprintf("VideoHub <%s>", username)
	// e.To = []string{to}
	// e.Subject = "Verification Code"
	// e.HTML = []byte(`
	// 	<h1>Verification Code</h1>
	// 	<p>Your verification code is: <strong>` + code + `</strong></p>`)
	// if err := e.Send(addr, smtp.PlainAuth("", username, password, host)); err != nil {
	// 	log.Println(err.Error()) // EOF
	// 	return err
	// }
	// return nil
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

func VerifyEmailCode(email string, code string) error {
	c, err := global.Rdb.Get(global.Ctx, email).Result()
	if err != nil {
		return errors.New("验证码已过期")
	}
	if c != code {
		return errors.New("验证码错误")
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
