package admin

type GetHistoricalDataRequest struct {
	StartDate string `form:"begin_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

type ListUsersRequest struct {
	Page   int    `form:"page" binding:"required"`
	Limit  int    `form:"limit" binding:"required"`
	Sort   string `form:"sort" binding:"required"`
	Order  *int8  `form:"order" binding:"required"`
	Status *int8  `form:"status" binding:"required"`
	Like   string `form:"like"`
	ID     uint   `form:"id"`
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
	Status *int   `json:"status"`
	Like   string `json:"like"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}
