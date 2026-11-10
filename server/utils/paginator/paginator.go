package paginator

import (
	"gorm.io/gorm"
)

func ParsePage(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size == 0 {
		return page, 0
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	return page, size
}

func PaginateQuery(query *gorm.DB, page, size int) *gorm.DB {
	page, size = ParsePage(page, size)
	if size > 0 {
		offset := (page - 1) * size
		query = query.Offset(offset).Limit(size)
	}
	return query
}

func Paginate[T any](query *gorm.DB, page, size int) ([]T, int64, error) {
	page, size = ParsePage(page, size)

	var count int64
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Model(new(T)).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var items []T
	if size > 0 {
		offset := (page - 1) * size
		query = query.Offset(offset).Limit(size)
	}
	if err := query.Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, count, nil
}

func HasMore(total int64, page, size int) bool {
	return int64(page*size) < total
}
