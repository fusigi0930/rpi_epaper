package disp

import(
	"os"
	"errors"
	"strconv"
	"strings"
	"image/color"

	"ggcal/log"
)

type GLargeDayControl struct {
	Ui *GObject

	Font *TTFFont
	Day *GTextControl
	Event *GTextControl

	DayFontPath string
	DayFontSize float64
	EventFontPath string
	EventFontSize float64
}

func (obj *GLargeDayControl) Init(parent GBase) {
	obj.Ui = &GObject{}
	obj.Ui.Init(parent)
	obj.Day = &GTextControl{}
	obj.Day.Init(obj)
	obj.Event = &GTextControl{}
	obj.Event.Init(obj)

	obj.AddChild(obj.Day)
	obj.AddChild(obj.Event)
}


func (obj *GLargeDayControl) SetPos(x int, y int) {
	obj.Ui.SetPos(x, y)
}

func (obj *GLargeDayControl) GetPos() (int, int) {
	return obj.Ui.GetPos()
}

func (obj *GLargeDayControl) GetAbsPos() (int, int) {
	return obj.Ui.GetAbsPos()
}

func (obj *GLargeDayControl) SetWidth(w int) {
	obj.Ui.SetWidth(w)
}

func (obj *GLargeDayControl) GetWidth() int {
	return obj.Ui.GetWidth()
}

func (obj *GLargeDayControl) SetHeight(h int) {
	obj.Ui.SetHeight(h)

}

func (obj *GLargeDayControl) GetHeight() int {
	return obj.Ui.GetHeight()
}

func (obj *GLargeDayControl) SetText(text string) {
	obj.Ui.SetText(text)
}

func (obj *GLargeDayControl) GetText() string {
	return obj.Ui.GetText()
}

func (obj *GLargeDayControl) SetTextColor(c color.Color) {
	obj.Ui.SetTextColor(c)
}

func (obj *GLargeDayControl) GetTextColor() color.Color {
	return obj.Ui.GetTextColor()
}

func (obj *GLargeDayControl) SetBorderColor(c color.Color) {
	obj.Ui.SetBorderColor(c)
}

func (obj *GLargeDayControl) GetBorderColor() color.Color {
	return obj.Ui.GetBorderColor()
}

func (obj *GLargeDayControl) SetBorder(b BorderType) {
	obj.Ui.SetBorder(b)
}

func (obj *GLargeDayControl) GetBorder() BorderType {
	return obj.Ui.GetBorder()
}

func (obj *GLargeDayControl) SetParent(parent GBase) {
	obj.Ui.SetParent(parent)
}

func (obj *GLargeDayControl) GetParent() GBase {
	return obj.Ui.GetParent()
}

func (obj *GLargeDayControl) Self() GBase {
	return obj
}

func (obj *GLargeDayControl) AddChild(child GBase) {
	obj.Ui.AddChild(child)
}

func (obj *GLargeDayControl) GetChild() []GBase {
	return obj.Ui.GetChild()
}

func (obj *GLargeDayControl) GetSurfaceContext() *SurfaceContext {
	x, y := obj.GetAbsPos()
	w := obj.GetWidth()
	h := obj.GetHeight()
	sc := obj.Ui.SC.AddChild(x, y, w, h)
	return sc
}

func (obj *GLargeDayControl) SetSurfaceContext(sc *SurfaceContext) {
	obj.Ui.SC = sc
}

func (obj *GLargeDayControl) Draw(scs ...*SurfaceContext) {
	if obj.Ui.SC == nil {
		log.LogService().Errorf("no way....\n")
		return
	}

	text := obj.Ui.GetText()
	sc := obj.GetSurfaceContext()
	lines := strings.Split(text, "\n")
	x, y := obj.GetAbsPos()

	var bk_color, pan_color color.Color
	bk_color = color.White
	pan_color = obj.GetTextColor()
	
	sc.CleanToColor(bk_color)
	for i, t := range lines {
		sc.DrawText(x+1, y+int(obj.Font.Size+1)*(i+1), t, obj.Font, pan_color)
	}

	obj.Ui.Draw()
}

func (obj *GLargeDayControl) Close() {
	obj.Ui.Close()
}

func (obj *GLargeDayControl) DumpParam(lead string) {
	obj.Ui.DumpParam(lead+"  ")
}

func (obj *GLargeDayControl) UpdateParam(params map[string] string) {
	obj.Ui.UpdateParam(params)

	fs, ok := params["fontsize"]
	var size float64
	if ok {
		size, _ = strconv.ParseFloat(fs, 64)
	}

	ft, ok := params["font"]
	if ok {
		obj.SetFont(ft, size)
	}

	df, ok := params["day_font"]
	if ok {
		obj.DayFontPath = df
	}
	dfs, ok := params["day_fontsize"]
	if ok {
		obj.DayFontSize, _ = strconv.ParseFloat(dfs, 64)
		obj.Day.SetFont(obj.DayFontPath, obj.DayFontSize)
	}

	ef, ok := params["event_font"]
	if ok {
		obj.EventFontPath = ef
	}
	efs, ok := params["event_fontsize"]
	if ok {
		obj.EventFontSize, _ = strconv.ParseFloat(efs, 64)
		obj.Event.SetFont(obj.EventFontPath, obj.EventFontSize)
	}

	eventx, _ := strconv.Atoi(params["event_x"])
	eventy, _ := strconv.Atoi(params["event_y"])
	obj.Event.SetPos(eventx, eventy)
	obj.Event.SetHeight(obj.GetHeight() - int(obj.Event.Font.Size) - 5)
	obj.Event.SetWidth(obj.GetWidth()-2)

	//textsize := obj.StringPixel() + 2
	obj.Day.SetPos(50, 2)
	obj.Day.SetHeight(int(obj.Day.Font.Size)+3)
	obj.Day.SetWidth(obj.GetWidth()-2-50)

	uniid, ok := params["id"]
	if ok {
		SetUniqueId(uniid, obj)
	}
}

func (obj *GLargeDayControl) SetDayText(t string) {
	obj.Day.SetText(t)
}

func (obj *GLargeDayControl) SetDayTextColor(c color.Color) {
	obj.Day.SetTextColor(c)
}

func (obj *GLargeDayControl) SetDayTextFont(file string, size float64) {
	obj.Day.SetFont(file, size)
}

func (obj *GLargeDayControl) SetEventText(t string) {
	obj.Event.SetText(t)
}

func (obj *GLargeDayControl) SetEventTextColor(c color.Color) {
	obj.Event.SetTextColor(c)
}

func (obj *GLargeDayControl) SetEventTextFont(file string, size float64) {
	obj.Event.SetFont(file, size)
}

func (obj *GLargeDayControl) SetFont(file string, size float64) {
	_, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) {
		log.LogService().Errorf("the file %s is not exist\n", file)
		return
	}

	obj.Font, _ = obj.Ui.SC.LoadFont(file, size, 72.0)
}

func (obj *GLargeDayControl) StringPixel() int {
	pxl, _ := obj.GetSurfaceContext().StringPixel(obj.GetText(), obj.Font)
	return pxl
}