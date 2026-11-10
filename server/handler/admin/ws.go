package admin

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"time"

	"server/database/connection"
	"server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type wsStats struct {
	Timestamp     string  `json:"timestamp"`
	Goroutines    int     `json:"goroutines"`
	MemAllocMB    float64 `json:"memAllocMB"`
	HeapAllocMB   float64 `json:"heapAllocMB"`
	GCPauseMs     float64 `json:"gcPauseMs"`
	NumGC         uint32  `json:"numGC"`
	OpenConns     int     `json:"openConns"`
	IdleConns     int     `json:"idleConns"`
	InUse         int     `json:"inUse"`
	TotalRequests int64   `json:"totalRequests"`
	Uptime        string  `json:"uptime"`
}

func (h *AdminHandler) StatsWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket 升级失败: %v", err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	for {
		select {
		case <-ticker.C:
			stats := collectWSStats(h.startTime)
			data, _ := json.Marshal(stats)
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		}
	}
}

func collectWSStats(startTime time.Time) *wsStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	var openConns, idleConns, inUse int
	if connection.DB != nil {
		if sqlDB, err := connection.DB.DB(); err == nil {
			poolStats := sqlDB.Stats()
			openConns = poolStats.OpenConnections
			idleConns = poolStats.Idle
			inUse = poolStats.InUse
		}
	}

	var gcPause float64
	if m.NumGC > 0 {
		gcPause = float64(m.PauseTotalNs) / float64(m.NumGC) / 1e6
	}

	return &wsStats{
		Timestamp:     time.Now().Format("15:04:05"),
		Goroutines:    runtime.NumGoroutine(),
		MemAllocMB:    float64(m.Alloc) / 1024 / 1024,
		HeapAllocMB:   float64(m.HeapAlloc) / 1024 / 1024,
		GCPauseMs:     gcPause,
		NumGC:         m.NumGC,
		OpenConns:     openConns,
		IdleConns:     idleConns,
		InUse:         inUse,
		TotalRequests: middleware.GetTotalRequests(),
		Uptime:        time.Since(startTime).Truncate(time.Second).String(),
	}
}
