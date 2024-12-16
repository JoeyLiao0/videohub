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
