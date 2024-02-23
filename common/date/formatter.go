package date

import (
	"log"
	"time"
)

const DATE_LAYOUT = "2006-01-02T15:04:05.999Z"

func FormatDate(date string) string {
	parsedDate, err := time.Parse(DATE_LAYOUT, date)
	if err != nil {
		log.Default().Printf("Error parsing time: %v", err)
	}
	return parsedDate.Format("2006.01.02")
}
