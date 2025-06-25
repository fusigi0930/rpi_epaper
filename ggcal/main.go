//go:build !disp_test && !cal_test

package main

import(
	"fmt"
	"time"

	"ggcal/log"
	"ggcal/cal"
	"ggcal/disp"
)

func updateTitle() {
	mon := disp.GetControlById("titleMonth")
	if mon != nil {
		mon.SetText(fmt.Sprintf("%s月", time.Now().Format("1")))
		ctrl := mon.(*disp.GTextControl)
		textsize := ctrl.StringPixel() + 5
		w := ctrl.GetParent().GetWidth()
		_, y := ctrl.GetPos()
		ctrl.SetWidth(textsize)
		ctrl.SetPos(w-textsize, y)
	}
	year := disp.GetControlById("titleYear")
	if year != nil {
		year.SetText(fmt.Sprintf("%s年", time.Now().Format("2006")))
		ctrl := year.(*disp.GTextControl)
		textsize := ctrl.StringPixel() + 2
		x, _ := mon.GetPos()
		_, y := ctrl.GetPos()
		ctrl.SetWidth(textsize)
		ctrl.SetPos(x-textsize, y)
	}
}

func updateLargeDays() {
	today := time.Now()
	ids := []string{"eventToday", "eventTomorrow", "eventAfter"}
	weekday := []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}

	eventlist := cal.GetDaysEvents(cal.CalService(), today.Format("2006-01-02"), 3)

	for i, id := range ids {
		check_day := today.Add(time.Duration(i) * 24 * time.Hour)
		log.LogService().Debugf("check date: %s\n", check_day.Format("2006-01-02"))
		wd := cal.GetWeekDay(check_day.Format("2006-01-02"))
		eventday := disp.GetControlById(id)
		if eventday != nil {
			ctrl := eventday.(*disp.GLargeDayControl)
			ctrl.SetDayText(fmt.Sprintf("%s(%s)", check_day.Format("01/02"), weekday[wd]))
			e, ok := eventlist[check_day.Format("2006-01-02")]
			if ok {
				ctrl.SetEventText(e)
			} else {
				ctrl.SetEventText("<今天沒事>")
			}
		}
	}
}

func updateCalendar() {
	grid := disp.GetControlById("calGrid")
	if grid != nil {
		ctrl := grid.(*disp.GCalGridControl)
		ctrl.UpdateGrid()
	}
}

func main() {
	log.InitLog("ggcal.log")
	log.LogService().Printf("initial log service...\n")

	err := disp.InitWin(1304, 984)
	if err != nil {
		log.LogService().Errorf("init virtual windows failed: %v\n", err)
	}

	go func() {
		if err := disp.LoadDef("layout.yaml", nil); err != nil {
			log.LogService().Errorf("no!!!! %v\n", err)
		}

		updateCalendar()
		updateTitle()
		updateLargeDays()

		disp.GetRootScreen().Draw()

	}()
	disp.EventLoop()

	disp.CloseWin()
}