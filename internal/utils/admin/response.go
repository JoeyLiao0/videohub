package admin

type GetRealTimeDataResponse struct {
	CPUPercent float64 `json:"cpu_percent"`
	MemTotal   float64 `json:"memory_total"`
	MemUsed    float64 `json:"memory_used"`
	OnlineNum  int     `json:"online_num"`
}

type Item struct {
	Date  string `json:"date"`
	Value int    `json:"value"`
}

type GetHistoricalDataResponse struct {
	Line []Item `json:"line"`
	Area []Item `json:"area"`
	Bar  []Item `json:"bar"`
}

type UserInfo struct {
	ID        uint   `json:"id" binding:"required"`
	Username  string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Avatar    string `json:"avatar" binding:"required"`
	Status    int8   `json:"status" binding:"required"`
	CreatedAt int64  `json:"time" binding:"required"`
}

type PageInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

type ListUsersResponse struct {
	Users []UserInfo `json:"users"`
	Pages PageInfo   `json:"pages"`
}
