package admin

type GetHistoricalDataRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

type ListUsersRequest struct {
	Page   int    `form:"page" binding:"required"`   // 当前页码，从1开始
	Limit  int    `form:"limit" binding:"required"`  // 每页数量
	Sort   string `form:"sort" binding:"required"`   // 排序字段，如id
	Order  *int8  `form:"order" binding:"required"`  // 排序方式：0-升序，1-降序
	Status *int8  `form:"status" binding:"required"` // 用户状态：-1全部，0正常，1禁用，2注销
	Like   string `form:"like"`                      // 用户名模糊匹配
	ID     uint   `form:"id"`                        // 用户ID，用于精确查询
}

type CreateUserRequest struct {
	Username string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required" validate:"email"`
	Avatar   string `json:"avatar"`
}

type UpdateUserRequest struct {
	ID     uint  `json:"uid" binding:"required"`
	Status *int8 `json:"new_status" binding:"required"`
}

type GetVideosRequest struct {
	Status *int   `json:"status"` // -1-全部 0-正常 1-审核 2-审核未通过 3-封禁
	Like   string `json:"like"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}
