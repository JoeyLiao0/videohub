package service

import (
	"math"
	"net/http"
	"videohub/internal/repository"
	"videohub/internal/utils"
	"videohub/internal/utils/admin"

	"github.com/sirupsen/logrus"
)

type UserList struct {
	//用户列表服务，只用到user表操作，所以只用注入user_repostiory
	userRepo *repository.User
}

// 工厂函数，返回单例的服务层操作对象
func NewUserList(ur *repository.User) *UserList {
	return &(UserList{userRepo: ur})
}

// GetUsers 获取用户列表
func (ul *UserList) GetUsers(request *admin.ListUsersRequest) *utils.Response {
	// 构建查询条件
	conditions := make(map[string]interface{})
	conditions["role"] = 0
	conditions["status"] = *request.Status
	if request.ID > 0 {
		conditions["id"] = request.ID
	}

	// 计算总记录数
	total, err := ul.userRepo.Count(conditions)
	if err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}
	totalPages := int(math.Ceil(float64(total) / float64(request.Limit)))

	fields := []string{"id", "username", "email", "avatar", "status", "created_at"}

	orderStr := request.Sort
	if *request.Order == 1 {
		orderStr += " DESC"
	} else {
		orderStr += " ASC"
	}

	offset := (request.Page - 1) * request.Limit

	var users []admin.UserInfo
	if err := ul.userRepo.SearchWithOrder(conditions, offset, request.Limit, orderStr, request.Like, fields, &users); err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	response := &admin.ListUsersResponse{
		Users: users,
		Pages: admin.PageInfo{
			Page:       request.Page,
			Limit:      request.Limit,
			TotalPages: totalPages,
		},
	}

	return utils.Ok(http.StatusOK, response)
}
