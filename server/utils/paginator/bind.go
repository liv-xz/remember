package paginator

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func BindPage(c *gin.Context) (int, int) {
	page, size := 1, 10
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			page = v
		}
	}
	if s := c.Query("size"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			size = v
		}
	}
	return ParsePage(page, size)
}
