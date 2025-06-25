//go:build cal_test

package main

import(
	"fmt"
	"time"

	"ggcal/log"
	"ggcal/cal"
)

func main() {
	log.InitLog("ggcal.log")
	log.LogService().Printf("initial log service...\n")

	if err := cal.LoadDef("calconfig.yaml"); err != nil {
		log.LogService().Errorf("no!!!!!!!: %v\n", err)
		return
	}
	service := cal.CalService()
	if service == nil {
		log.LogService().Errorf("no!!!!!!! \n")
		return
	}
	//b, s := cal.IsHoliday(service, "2025-01-01")
	//fmt.Println(b, s)

	cal.GetEvents(service, "2025-01-01")

	fmt.Println(cal.GetWeekDay(time.Now().Format(time.RFC3339)))
	fmt.Println(cal.GetMonthDayCount("2100-02-02"))
	fmt.Println(cal.GetMonthDayCount("2024-02-02"))
	fmt.Println(cal.GetMonthDayCount("2025-02-02"))
	fmt.Println(cal.GetMonthDayCount("2026-02-02"))
	fmt.Println(cal.GetMonthDayCount("1582-10-02"))
}