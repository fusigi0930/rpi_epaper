package disp

import(
	"image/color"

	"ggcal/log"
	"ggcal/cal"
	"fmt"
	"time"
	"strconv"
)

type GCalGridControl struct {
	Ui *GObject
	Row	int
	Col int
	gridWidth int
	gridHeight int
	FirstIndex int

	DayFontPath string
	DayFontSize float64
	EventFontPath string
	EventFontSize float64
	WeekDayFontPath string
	WeekDayFontSize float64
	WeekDayFont *TTFFont
}

const OFFSET_HEAD int = 60

func (obj *GCalGridControl) Init(parent GBase) {
	obj.Ui = &GObject{}
	obj.Ui.Init(parent)

	for i := 0; i < 31; i++ {
		c := &GSmallDayCtrl{}
		c.Init(obj)
		c.SetDayText(fmt.Sprintf("%d", i+1))
		obj.AddChild(c)
	}
}

func (obj *GCalGridControl) SetPos(x int, y int) {
	obj.Ui.SetPos(x, y)
}

func (obj *GCalGridControl) GetPos() (int, int) {
	return obj.Ui.GetPos()
}

func (obj *GCalGridControl) GetAbsPos() (int, int) {
	return obj.Ui.GetAbsPos()
}

func (obj *GCalGridControl) SetWidth(w int) {
	obj.Ui.SetWidth(w)
}

func (obj *GCalGridControl) GetWidth() int {
	return obj.Ui.GetWidth()
}

func (obj *GCalGridControl) SetHeight(h int) {
	obj.Ui.SetHeight(h)
}

func (obj *GCalGridControl) GetHeight() int {
	return obj.Ui.GetHeight()
}

func (obj *GCalGridControl) SetText(text string) {
	obj.Ui.SetText(text)
}

func (obj *GCalGridControl) GetText() string {
	return obj.Ui.GetText()
}

func (obj *GCalGridControl) SetTextColor(c color.Color) {
	obj.Ui.SetTextColor(c)
}

func (obj *GCalGridControl) GetTextColor() color.Color {
	return obj.Ui.GetTextColor()
}

func (obj *GCalGridControl) SetBorderColor(c color.Color) {
	obj.Ui.SetBorderColor(c)
}

func (obj *GCalGridControl) GetBorderColor() color.Color {
	return obj.Ui.GetBorderColor()
}

func (obj *GCalGridControl) SetBorder(b BorderType) {
	obj.Ui.SetBorder(b)
}

func (obj *GCalGridControl) GetBorder() BorderType {
	return obj.Ui.GetBorder()
}

func (obj *GCalGridControl) SetParent(parent GBase) {
	obj.Ui.SetParent(parent)
}

func (obj *GCalGridControl) GetParent() GBase {
	return obj.Ui.GetParent()
}

func (obj *GCalGridControl) Self() GBase {
	return obj
}

func (obj *GCalGridControl) AddChild(child GBase) {
	obj.Ui.AddChild(child)
}

func (obj *GCalGridControl) GetChild() []GBase {
	return obj.Ui.GetChild()
}

func (obj *GCalGridControl) GetSurfaceContext() *SurfaceContext {
	x, y := obj.GetAbsPos()
	w := obj.GetWidth()
	h := obj.GetHeight()
	sc := obj.Ui.SC.AddChild(x, y, w, h)
	return sc
}

func (obj *GCalGridControl) SetSurfaceContext(sc *SurfaceContext) {
	obj.Ui.SetSurfaceContext(sc)
}

func (obj *GCalGridControl) Draw(scs ...*SurfaceContext) {

	if obj.Ui.SC == nil {
		log.LogService().Errorf("no way....\n")
		return
	}

	sc := obj.GetSurfaceContext()
	x, y := obj.GetAbsPos()
	grid_y := y+OFFSET_HEAD
	w := obj.GetWidth()
	h := obj.GetHeight()
	grid_h := h-OFFSET_HEAD
	offset_y := y + OFFSET_HEAD / 2

	sc.CleanToColor(color.White)
	wk := []string{"日", "一", "二", "三", "四", "五", "六"}

	for i := 1; i < obj.Col; i++ {
		week_x := x+obj.gridWidth*(i-1)+obj.gridWidth/2
		var c color.Color = color.Black
		if i == 1 {
			c = color.RGBA{R:255, A:255}
		}
		sc.DrawText(week_x, offset_y+int(obj.WeekDayFontSize), wk[i-1], obj.WeekDayFont, c)
		xc := x+(obj.gridWidth+1)*i
		sc.DrawDotLine(xc, grid_y, xc, grid_y+grid_h-1, color.Black)
	}
	week_x := x+obj.gridWidth*6+obj.gridWidth/2
	sc.DrawText(week_x, offset_y+int(obj.WeekDayFontSize), wk[6], obj.WeekDayFont, color.RGBA{R:255, A:255})

	for i := 1; i < obj.Row; i++ {
		yr := grid_y+(obj.gridHeight+1)*i
		sc.DrawDotLine(x, yr, x+w-1, yr, color.Black)
	}

	obj.Ui.Draw()

}

func (obj *GCalGridControl) Close() {
	obj.Ui.Close()
}

func (obj *GCalGridControl) DumpParam(lead string) {
	obj.Ui.DumpParam(lead+"  ")
}

func (obj *GCalGridControl) UpdateParam(params map[string] string) {
	obj.Ui.UpdateParam(params)

	wdf, ok := params["weekday_font"]
	if ok {
		obj.WeekDayFontPath = wdf
	}
	wdfs, ok := params["weekday_fontsize"]
	if ok {
		obj.WeekDayFontSize, _ = strconv.ParseFloat(wdfs, 64)
		obj.WeekDayFont, _ = obj.Ui.SC.LoadFont(obj.WeekDayFontPath, obj.WeekDayFontSize, 72.0)
	}

	df, ok := params["day_font"]
	if ok {
		obj.DayFontPath = df
	}
	dfs, ok := params["day_fontsize"]
	if ok {
		obj.DayFontSize, _ = strconv.ParseFloat(dfs, 64)
	}

	ef, ok := params["event_font"]
	if ok {
		obj.EventFontPath = ef
	}
	efs, ok := params["event_fontsize"]
	if ok {
		obj.EventFontSize, _ = strconv.ParseFloat(efs, 64)
	}

	uniid, ok := params["id"]
	if ok {
		SetUniqueId(uniid, obj)
	}
}

func (obj *GCalGridControl) UpdateGrid() {
	today := time.Now()
	//today, _ := time.Parse("2006-01-02T15:04:05-07:00", "2025-01-01T18:45:00+08:00")
	first_weekday := cal.GetWeekDay(fmt.Sprintf("%s-01", today.Format("2006-01")))
	if first_weekday >= cal.Fri {
		obj.Row = 6
	} else {
		obj.Row = 5
	}
	obj.Col = 7
	obj.FirstIndex = int(first_weekday)
	log.LogService().Infof("first weekday in this month is: %d\n", obj.FirstIndex)

	w := obj.GetWidth()
	h := obj.GetHeight()

	obj.gridWidth = (w - 8) / obj.Col
	obj.gridHeight = (h - OFFSET_HEAD - obj.Row - 2) / (obj.Row)

	gx := obj.FirstIndex
	gy := 0
	dayCount := cal.GetMonthDayCount(today.Format("2006-01-02"))
	holidays := cal.GetHolidays(cal.CalService(), today.Format("2006-01-02"))
	eventlist := cal.GetEvents(cal.CalService(), today.Format("2006-01-02"))

	for i, ci := range obj.GetChild() {
		check_date := fmt.Sprintf("%s-%02d", today.Format("2006-01"), i+1)

		c := ci.(*GSmallDayCtrl)
		draw_x := 2 + (gx*(obj.gridWidth + 1))
		draw_y := OFFSET_HEAD + 2 + (gy*(obj.gridHeight + 1))
		c.SetDayTextFont(obj.DayFontPath, obj.DayFontSize)
		c.SetEventTextFont(obj.EventFontPath, obj.EventFontSize)
		c.SetHolidayTextFont(obj.EventFontPath, obj.EventFontSize)
		c.SetLunarTextFont(obj.EventFontPath, obj.EventFontSize)
		c.SetPos(draw_x, draw_y)
		c.Visible = (i < dayCount)
		if i >= dayCount {
			continue
		}

		holiday, ok := holidays[check_date]
		if ok {
			c.SetDayTextColor(color.RGBA{R:255, A:255})
			c.SetHolidayText(holiday)
		} else {
			c.SetHolidayText("")
		}

		ents, ok := eventlist[check_date]
		if ok {
			log.LogService().Debugf("set date: %s event %s\n", check_date, ents)
			c.SetEventText(ents)
		} else {
			c.SetEventText("")
		}

		if today.Format("2006-01-02") == check_date {
			c.Day.Inverse = true
		} else {
			c.Day.Inverse = false
		}

		chineseDate := cal.GetLunarDate(check_date)
		c.SetLunarText(chineseDate)

		if wd := cal.GetWeekDay(check_date); wd == 0 || wd == 6 {
			c.SetDayTextColor(color.RGBA{R:255, A:255})
		} else {
			c.SetDayTextColor(color.Black)
		}


		c.SetWidth(obj.gridWidth-2)
		c.SetHeight(obj.gridHeight-2)


		gx++
		if gx > 6 {
			gx = 0
			gy++
		}
	}
}