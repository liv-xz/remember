package admin

import (
	"net/http"
	"runtime"
	"strconv"
	"time"

	"server/middleware"
	oplogmodel "server/model/operation_log"
	"server/utils/paginator"
	"server/utils/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db        *gorm.DB
	startTime time.Time
}

func NewAdminHandler(db *gorm.DB, startTime time.Time) *AdminHandler {
	return &AdminHandler{db: db, startTime: startTime}
}

func (h *AdminHandler) StatsHandler(c *gin.Context) {
	stats := h.collectStats()
	response.SuccessData(c, stats)
}

func (h *AdminHandler) OverviewHandler(c *gin.Context) {
	var resp OverviewResponse
	h.db.Table("users").Where("status != 2").Count(&resp.TotalUsers)
	h.db.Table("graves").Count(&resp.TotalGraves)
	h.db.Table("deceased").Count(&resp.TotalDeceased)
	h.db.Table("memorial_records").Count(&resp.TotalRecords)
	h.db.Table("memorial_events").Count(&resp.TotalEvents)
	h.db.Table("grave_members").Count(&resp.TotalMembers)

	today := time.Now().Format("2006-01-02")
	h.db.Table("users").Where("DATE(created_at) = ? AND status != 2", today).Count(&resp.TodayUsers)
	h.db.Table("graves").Where("DATE(created_at) = ?", today).Count(&resp.TodayGraves)
	h.db.Table("deceased").Where("DATE(created_at) = ?", today).Count(&resp.TodayDeceased)
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	h.db.Table("users").Where("updated_at > ?", sevenDaysAgo).Count(&resp.ActiveUsers7d)

	response.SuccessData(c, &resp)
}

func (h *AdminHandler) TrendHandler(c *gin.Context) {
	days := 30
	if d, err := strconv.Atoi(c.DefaultQuery("days", "30")); err == nil && d > 0 {
		days = d
	}
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	type row struct {
		Date  string
		Count int64
	}

	var users []row
	h.db.Raw("SELECT DATE(created_at) as date, COUNT(*) as count FROM users WHERE created_at >= ? GROUP BY DATE(created_at) ORDER BY date", startDate).Scan(&users)
	userMap := make(map[string]int64)
	for _, r := range users {
		userMap[r.Date] = r.Count
	}

	var graves []row
	h.db.Raw("SELECT DATE(created_at) as date, COUNT(*) as count FROM graves WHERE created_at >= ? GROUP BY DATE(created_at) ORDER BY date", startDate).Scan(&graves)
	graveMap := make(map[string]int64)
	for _, r := range graves {
		graveMap[r.Date] = r.Count
	}

	var deceased []row
	h.db.Raw("SELECT DATE(created_at) as date, COUNT(*) as count FROM deceased WHERE created_at >= ? GROUP BY DATE(created_at) ORDER BY date", startDate).Scan(&deceased)
	deceasedMap := make(map[string]int64)
	for _, r := range deceased {
		deceasedMap[r.Date] = r.Count
	}

	items := make([]TrendItem, 0, days)
	for i := 0; i < days; i++ {
		d := time.Now().AddDate(0, 0, -(days - 1 - i)).Format("2006-01-02")
		items = append(items, TrendItem{
			Date:     d,
			Users:    userMap[d],
			Graves:   graveMap[d],
			Deceased: deceasedMap[d],
		})
	}

	response.SuccessData(c, &TrendResponse{Days: days, Items: items})
}

func (h *AdminHandler) HotHandler(c *gin.Context) {
	limit := 10
	if l, err := strconv.Atoi(c.DefaultQuery("limit", "10")); err == nil && l > 0 {
		limit = l
	}

	type graveRow struct {
		ID            uint
		Name          string
		CreatedAt     time.Time
		DeceasedCount int64
		MemberCount   int64
		RecordCount   int64
	}

	var rows []graveRow
	h.db.Raw(`
		SELECT g.id, g.name, g.created_at,
			COALESCE(dc.cnt, 0)   AS deceased_count,
			COALESCE(mc.cnt, 0)   AS member_count,
			COALESCE(rc.cnt, 0)   AS record_count
		FROM graves g
		LEFT JOIN (SELECT grave_id, COUNT(*) AS cnt FROM deceased GROUP BY grave_id) dc ON dc.grave_id = g.id
		LEFT JOIN (SELECT grave_id, COUNT(*) AS cnt FROM grave_members GROUP BY grave_id) mc ON mc.grave_id = g.id
		LEFT JOIN (SELECT grave_id, COUNT(*) AS cnt FROM memorial_records GROUP BY grave_id) rc ON rc.grave_id = g.id
		ORDER BY (COALESCE(dc.cnt, 0) + COALESCE(mc.cnt, 0) + COALESCE(rc.cnt, 0)) DESC
		LIMIT ?
	`, limit).Scan(&rows)

	items := make([]HotGraveItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, HotGraveItem{
			ID:            r.ID,
			Name:          r.Name,
			DeceasedCount: r.DeceasedCount,
			MemberCount:   r.MemberCount,
			RecordCount:   r.RecordCount,
			CreatedAt:     r.CreatedAt.Format("2006-01-02"),
		})
	}
	response.SuccessData(c, items)
}

func (h *AdminHandler) ListUsersHandler(c *gin.Context) {
	page, size := paginator.BindPage(c)
	search := c.Query("search")

	query := h.db.Model(&userRow{}).Where("status != 2")
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?", like, like, like)
	}
	query = query.Order("id DESC")

	items, total, err := paginator.Paginate[userRow](query, page, size)
	if err != nil {
		response.Fail(c, 500)
		return
	}

	list := make([]UserListItem, 0, len(items))
	for _, u := range items {
		list = append(list, UserListItem{
			ID:        u.ID,
			Username:  u.Username,
			Nickname:  u.Nickname,
			Email:     u.Email,
			Avatar:    u.Avatar,
			Status:    u.Status,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04"),
		})
	}
	response.SuccessPage(c, list, total, page, size)
}

func (h *AdminHandler) ListGravesHandler(c *gin.Context) {
	page, size := paginator.BindPage(c)

	type row struct {
		ID            uint
		Name          string
		GraveType     int8
		GraveMode     int8
		Status        int
		CreatedBy     uint
		CreatorName   string
		DeceasedCount int64 `gorm:"->"`
		MemberCount   int64 `gorm:"->"`
		CreatedAt     time.Time
	}

	sub := h.db.Table("graves").Select("graves.*, u.nickname as creator_name")
	sub = sub.Joins("LEFT JOIN users u ON u.id = graves.created_by")
	query := h.db.Table("(?) AS g", sub).
		Select("g.*, COALESCE(dc.cnt,0) AS deceased_count, COALESCE(mc.cnt,0) AS member_count").
		Joins("LEFT JOIN (SELECT grave_id, COUNT(*) AS cnt FROM deceased GROUP BY grave_id) dc ON dc.grave_id = g.id").
		Joins("LEFT JOIN (SELECT grave_id, COUNT(*) AS cnt FROM grave_members GROUP BY grave_id) mc ON mc.grave_id = g.id").
		Order("g.id DESC")

	items, total, err := paginator.Paginate[row](query, page, size)
	if err != nil {
		response.Fail(c, 500)
		return
	}

	list := make([]GraveListItem, 0, len(items))
	for _, g := range items {
		list = append(list, GraveListItem{
			ID:            g.ID,
			Name:          g.Name,
			GraveType:     g.GraveType,
			GraveMode:     g.GraveMode,
			Status:        g.Status,
			Creator:       g.CreatorName,
			DeceasedCount: g.DeceasedCount,
			MemberCount:   g.MemberCount,
			CreatedAt:     g.CreatedAt.Format("2006-01-02 15:04"),
		})
	}
	response.SuccessPage(c, list, total, page, size)
}

func (h *AdminHandler) UpdateUserStatusHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest)
		return
	}

	var req UpdateUserStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest)
		return
	}

	if req.Status < 0 || req.Status > 1 {
		response.Fail(c, http.StatusBadRequest)
		return
	}

	result := h.db.Model(&userRow{}).Where("id = ?", id).Update("status", req.Status)
	if result.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound)
		return
	}
	response.Success(c)
}

func (h *AdminHandler) DeleteUserHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest)
		return
	}

	result := h.db.Model(&userRow{}).Where("id = ?", id).Update("status", 2)
	if result.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound)
		return
	}
	response.Success(c)
}

type userRow struct {
	ID        uint
	Username  string
	Nickname  string
	Email     string
	Avatar    string
	Status    int8
	CreatedAt time.Time
}

func (userRow) TableName() string { return "users" }

func (h *AdminHandler) collectStats() *StatsResponse {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	sqlDB, err := h.db.DB()
	dbStats := DBStats{}
	if err == nil {
		poolStats := sqlDB.Stats()
		dbStats.OpenConns = poolStats.OpenConnections
		dbStats.IdleConns = poolStats.Idle
		dbStats.MaxOpenConns = poolStats.MaxOpenConnections
	}

	h.db.Table("users").Where("status != 2").Count(&dbStats.Users)
	h.db.Table("graves").Count(&dbStats.Graves)
	h.db.Table("deceased").Count(&dbStats.Deceased)
	h.db.Table("grave_members").Count(&dbStats.GraveMembers)
	h.db.Table("memorial_records").Count(&dbStats.MemorialRecords)
	h.db.Table("memorial_events").Count(&dbStats.MemorialEvents)
	h.db.Table("kinship_relations").Count(&dbStats.KinshipRelations)

	var ops []oplogmodel.OperationLog
	h.db.Order("created_at DESC").Limit(20).Find(&ops)

	recentOps := make([]OpLogItem, 0, len(ops))
	for _, op := range ops {
		recentOps = append(recentOps, OpLogItem{
			Time:   op.CreatedAt.Format("15:04:05"),
			Action: op.Action,
			Target: op.TargetType,
			Status: op.ResponseStatus,
		})
	}

	return &StatsResponse{
		Server: ServerStats{
			Uptime:        time.Since(h.startTime).Truncate(time.Second).String(),
			GoVersion:     runtime.Version(),
			Goroutines:    runtime.NumGoroutine(),
			MemAllocMB:    float64(m.Alloc) / 1024 / 1024,
			TotalRequests: middleware.GetTotalRequests(),
		},
		Database:  dbStats,
		RecentOps: recentOps,
	}
}
