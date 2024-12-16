package service

import (
	"net/http"
	"time"
	"videohub/global"
	"videohub/internal/utils"
	"videohub/internal/utils/admin"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/sirupsen/logrus"
)

type Data struct {
}

func NewData() *Data {
	return &Data{}
}

func (e *Data) GetRealTimeData() *utils.Response {
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	virtualMem, err := mem.VirtualMemory()
	if err != nil {
		logrus.Error(err.Error())
		return utils.Error(http.StatusInternalServerError, "服务器内部错误")
	}

	pattern := "user:*:is_online"
	var cursor uint64
	var count int
	for {
		keys, nextCursor, err := global.Rdb.Scan(global.Ctx, cursor, pattern, 100).Result()
		if err != nil {
			return utils.Error(http.StatusInternalServerError, "服务器内部错误")
		}
		count += len(keys)
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	return utils.Ok(http.StatusOK, &admin.GetRealTimeDataResponse{
		CPUPercent: cpuPercent[0],
		MemTotal:   float64(virtualMem.Total) / 1e9,
		MemUsed:    float64(virtualMem.Used) / 1e9,
		OnlineNum:  count,
	})
}

func (e *Data) GetHistoricalData() *utils.Response {
	// TODO
	return utils.Ok(http.StatusOK, nil)
}
