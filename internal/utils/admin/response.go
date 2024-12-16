package admin

type GetRealTimeDataResponse struct {
	CPUPercent float64 `json:"cpu_percent"`
	MemTotal   float64 `json:"memory_total"`
	MemUsed    float64 `json:"memory_used"`
	OnlineNum  int     `json:"online_num"`
}

type GetHistoricalDataResponse struct {
}
