package utils

import (
	"time"

	"github.com/gin-gonic/gin"
)

func BuildFilterFromQuery(ctx *gin.Context) map[string]interface{} {
	filter := make(map[string]interface{})
	dateKeys := map[string]bool{
		"start_date":        true,
		"end_date":          true,
		"budget_start_date": true,
		"budget_end_date":   true,
	}
	for key, values := range ctx.Request.URL.Query() {
		if len(values) > 0 && values[0] != "" {
			if dateKeys[key] {
				parsedDate, err := time.Parse("2006-01-02", values[0])
				if err == nil {
					filter[key] = parsedDate
					continue
				}
			}
			filter[key] = values[0]
		}
	}
	return filter
}
