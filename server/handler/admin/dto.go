package admin

type StatsResponse struct {
	Server    ServerStats `json:"server"`
	Database  DBStats     `json:"database"`
	RecentOps []OpLogItem `json:"recentOps"`
}

type ServerStats struct {
	Uptime        string  `json:"uptime"`
	GoVersion     string  `json:"goVersion"`
	Goroutines    int     `json:"goroutines"`
	MemAllocMB    float64 `json:"memAllocMB"`
	TotalRequests int64   `json:"totalRequests"`
}

type DBStats struct {
	OpenConns        int   `json:"openConns"`
	IdleConns        int   `json:"idleConns"`
	MaxOpenConns     int   `json:"maxOpenConns"`
	Users            int64 `json:"users"`
	Graves           int64 `json:"graves"`
	Deceased         int64 `json:"deceased"`
	GraveMembers     int64 `json:"graveMembers"`
	MemorialRecords  int64 `json:"memorialRecords"`
	MemorialEvents   int64 `json:"memorialEvents"`
	KinshipRelations int64 `json:"kinshipRelations"`
}

type OpLogItem struct {
	Time   string `json:"time"`
	Action string `json:"action"`
	Target string `json:"target"`
	Status int    `json:"status"`
}

type OverviewResponse struct {
	TotalUsers    int64 `json:"totalUsers"`
	TotalGraves   int64 `json:"totalGraves"`
	TotalDeceased int64 `json:"totalDeceased"`
	TotalRecords  int64 `json:"totalRecords"`
	TotalEvents   int64 `json:"totalEvents"`
	TodayUsers    int64 `json:"todayUsers"`
	TodayGraves   int64 `json:"todayGraves"`
	TodayDeceased int64 `json:"todayDeceased"`
	TotalMembers  int64 `json:"totalMembers"`
	ActiveUsers7d int64 `json:"activeUsers7d"`
}

type TrendItem struct {
	Date     string `json:"date"`
	Users    int64  `json:"users"`
	Graves   int64  `json:"graves"`
	Deceased int64  `json:"deceased"`
}

type TrendResponse struct {
	Days  int         `json:"days"`
	Items []TrendItem `json:"items"`
}

type HotGraveItem struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	DeceasedCount int64  `json:"deceasedCount"`
	MemberCount   int64  `json:"memberCount"`
	RecordCount   int64  `json:"recordCount"`
	CreatedAt     string `json:"createdAt"`
}

type UserListItem struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	Status    int8   `json:"status"`
	CreatedAt string `json:"createdAt"`
}

type GraveListItem struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	GraveType     int8   `json:"graveType"`
	GraveMode     int8   `json:"graveMode"`
	Status        int    `json:"status"`
	Creator       string `json:"creator"`
	DeceasedCount int64  `json:"deceasedCount"`
	MemberCount   int64  `json:"memberCount"`
	CreatedAt     string `json:"createdAt"`
}

type UpdateUserStatusReq struct {
	Status int8 `json:"status" binding:"required"`
}
