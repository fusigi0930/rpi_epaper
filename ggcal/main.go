//go:build !disp_test && !cal_test

package main

import(
	"fmt"
	"time"
	"os"
	"sync"

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

func getNextDuration() time.Duration {
	targetHours := []int{0, 6, 9, 12, 15, 18, 21}
	now := time.Now()
	var nextRun time.Time

	for _, hour := range targetHours {
		potentialNextRun := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())
		if potentialNextRun.After(now) {
			nextRun = potentialNextRun
			break
		}
	}

	if nextRun.IsZero() {
		tomorrow := now.Add(24 * time.Hour)
		nextRun = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), targetHours[0], 0, 0, 0, now.Location())
	}

	log.LogService().Infof("next update: %s will run %s\n", nextRun.Sub(now).Round(time.Second), nextRun.Format("2006-01-02 15:04:05"))
	return nextRun.Sub(now)
}

func runSchedule(drawSignal chan<- struct{}) {
	var wg sync.WaitGroup
	for {
		duration := getNextDuration()
		timer := time.NewTimer(duration)
		<-timer.C
		wg.Add(3)
		go func() {
			updateCalendar()
			wg.Done()
		}()
		go func() {
			updateTitle()
			wg.Done()
		}()
		go func() {
			updateLargeDays()
			wg.Done()
		}()
		wg.Wait()
		drawSignal <- struct{}{}
	}
}


func main() {
	log.InitLog("ggcal.log")
	log.LogService().Printf("initial log service...\n")

	var mWg sync.WaitGroup

	mWg.Add(1)
	go func() {
		if err := disp.LoadDef("layout.yaml", nil); err != nil {
			log.LogService().Errorf("no!!!! %v\n", err)
		}
		mWg.Done()
	}()

	log.LogService().Infof("initial display windows in the other thread\n")
	err := disp.InitWin(1304, 984)
	if err != nil {
		log.LogService().Errorf("init virtual windows failed: %v\n", err)
		os.Exit(1)
	}

	mWg.Wait()

	mWg.Add(3)
	go func() {
		updateCalendar()
		mWg.Done()
	}()
	go func() {
		updateTitle()
		mWg.Done()
	}()
	go func() {
		updateLargeDays()
		mWg.Done()
	}()
	
	mWg.Wait()
	log.LogService().Infof("update window\n")
	disp.GetRootScreen().Draw()

	drawSignal := make(chan struct{}, 1)
	go runSchedule(drawSignal)

	disp.EventLoop(drawSignal)

	log.LogService().Infof("ready for closing window\n")
	disp.CloseWin()
}