package toolbox

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PaginationResult struct {
	Data      interface{}
	Page      int
	PageSize  int
	TotalRows int64
}

func Paginate(c *fiber.Ctx, db *gorm.DB, model interface{}, preloads []string, searchFields []string) (*PaginationResult, error) {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	search := c.Query("search", "")
	offset := (page - 1) * pageSize

	query := db.Model(model)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if search != "" {
		for i, field := range searchFields {
			if i == 0 {
				query = query.Where(field+" LIKE ?", "%"+search+"%")
			} else {
				query = query.Or(field+" LIKE ?", "%"+search+"%")
			}
		}
	}

	var totalRows int64
	db.Model(model).Count(&totalRows)
	query.Offset(offset).Limit(pageSize).Find(model)

	return &PaginationResult{
		Data: model, Page: page, PageSize: pageSize, TotalRows: totalRows,
	}, nil
}
