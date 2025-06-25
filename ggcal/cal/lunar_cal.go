package cal

import (
	"fmt"
	"strings"
	"strconv"

	"ggcal/log"

	"github.com/Lofanmi/chinese-calendar-golang/calendar"
)


func GetLunarDate(date string) string {
	ds := strings.Split(date, "-")
	if len(ds) < 3 {
		log.LogService().Errorf("date format error\n")
		return ""
	}
	y, _ := strconv.Atoi(ds[0])
	m, _ := strconv.Atoi(ds[1])
	d, _ := strconv.Atoi(ds[2])
	chineseDate := calendar.BySolar(int64(y), int64(m), int64(d), int64(10), int64(10), int64(10))
	retDate := fmt.Sprintf("%s%s", chineseDate.Lunar.MonthAlias(), chineseDate.Lunar.DayAlias())
	return retDate
}