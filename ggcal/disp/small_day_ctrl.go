package disp

import(
	"strings"
	"image/color"

	"ggcal/log"
)

type GSmallDayCtrl struct {
	Ui *GObject

	Day *GTextControl
	Event *GTextControl
	Holiday *GTextControl
	Lunar *GTextControl

	Visible bool
}

func (obj *GSmallDayCtrl) Init(parent GBase) {
	obj.Ui = &GObject{}
	obj.Ui.Init(parent)
	obj.Day = &GTextControl{}
	obj.Day.Init(obj)
	obj.Event = &GTextControl{}
	obj.Event.Init(obj)
	obj.Holiday = &GTextControl{}
	obj.Holiday.Init(obj)
	obj.Lunar = &GTextControl{}
	obj.Lunar.Init(obj)


	obj.AddChild(obj.Day)
	obj.AddChild(obj.Event)
	obj.AddChild(obj.Holiday)
	obj.AddChild(obj.Lunar)

	//obj.Event.SetText("10:20 xxxxxxxxxxx\n14:33 測試啊~\n龍龜賽跑")
}

func (obj *GSmallDayCtrl) SetPos(x int, y int) {
	obj.Ui.SetPos(x, y)
}

func (obj *GSmallDayCtrl) GetPos() (int, int) {
	return obj.Ui.GetPos()
}

func (obj *GSmallDayCtrl) GetAbsPos() (int, int) {
	return obj.Ui.GetAbsPos()
}

func (obj *GSmallDayCtrl) SetWidth(w int) {
	obj.Ui.SetWidth(w)
	textsize := obj.Day.StringPixel() + 5
	obj.Day.SetPos(w - textsize, 5)
	obj.Day.SetWidth(textsize+2)
	obj.Event.SetWidth(w-10)
	obj.Lunar.SetPos(5, 5)
	obj.Lunar.SetWidth(w-5-(textsize+2))
	obj.Holiday.SetPos(5, 5 + int(obj.Lunar.Font.Size + 3))
	obj.Holiday.SetWidth(w-5-(textsize+2))
}

func (obj *GSmallDayCtrl) GetWidth() int {
	return obj.Ui.GetWidth()
}

func (obj *GSmallDayCtrl) SetHeight(h int) {
	obj.Ui.SetHeight(h)
	textsize := int(obj.Day.Font.Size + 2)
	obj.Day.SetHeight(textsize+4)
	line_count := len(strings.Split(obj.Event.GetText(), "\n"))
	if line_count > 0 {
		obj.Event.SetPos(5, h - int(obj.Event.Font.Size+3) * line_count)
		obj.Event.SetHeight(int(obj.Event.Font.Size+3) * line_count)
	}
	obj.Lunar.SetHeight(int(obj.Lunar.Font.Size + 3))
	obj.Holiday.SetHeight(int(obj.Holiday.Font.Size + 3))
}

func (obj *GSmallDayCtrl) GetHeight() int {
	return obj.Ui.GetHeight()
}

func (obj *GSmallDayCtrl) SetText(text string) {
	obj.Ui.SetText(text)
}

func (obj *GSmallDayCtrl) GetText() string {
	return obj.Ui.GetText()
}

func (obj *GSmallDayCtrl) SetTextColor(c color.Color) {
	obj.Ui.SetTextColor(c)
}

func (obj *GSmallDayCtrl) GetTextColor() color.Color {
	return obj.Ui.GetTextColor()
}

func (obj *GSmallDayCtrl) SetBorderColor(c color.Color) {
	obj.Ui.SetBorderColor(c)
}

func (obj *GSmallDayCtrl) GetBorderColor() color.Color {
	return obj.Ui.GetBorderColor()
}

func (obj *GSmallDayCtrl) SetBorder(b BorderType) {
	obj.Ui.SetBorder(b)
}

func (obj *GSmallDayCtrl) GetBorder() BorderType {
	return obj.Ui.GetBorder()
}

func (obj *GSmallDayCtrl) SetParent(parent GBase) {
	obj.Ui.SetParent(parent)
}

func (obj *GSmallDayCtrl) GetParent() GBase {
	return obj.Ui.GetParent()
}

func (obj *GSmallDayCtrl) Self() GBase {
	return obj
}

func (obj *GSmallDayCtrl) AddChild(child GBase) {
	obj.Ui.AddChild(child)
}

func (obj *GSmallDayCtrl) GetChild() []GBase {
	return obj.Ui.GetChild()
}

func (obj *GSmallDayCtrl) GetSurfaceContext() *SurfaceContext {
	x, y := obj.GetAbsPos()
	w := obj.GetWidth()
	h := obj.GetHeight()
	sc := obj.Ui.SC.AddChild(x, y, w, h)
	return sc
}

func (obj *GSmallDayCtrl) SetSurfaceContext(sc *SurfaceContext) {
	obj.Ui.SC = sc
}

func (obj *GSmallDayCtrl) Draw(scs ...*SurfaceContext) {
	if obj.Ui.SC == nil {
		log.LogService().Errorf("no way....\n")
		return
	}
	if obj.Visible == false {
		return
	}

	obj.Ui.Draw()
}

func (obj *GSmallDayCtrl) Close() {
	obj.Ui.Close()
}

func (obj *GSmallDayCtrl) DumpParam(lead string) {
	obj.Ui.DumpParam(lead+"  ")
}

func (obj *GSmallDayCtrl) UpdateParam(params map[string] string) {
	obj.Ui.UpdateParam(params)

}

func (obj *GSmallDayCtrl) SetDayText(t string) {
	obj.Day.SetText(t)
}

func (obj *GSmallDayCtrl) SetDayTextColor(c color.Color) {
	obj.Day.SetTextColor(c)
}

func (obj *GSmallDayCtrl) SetDayTextFont(file string, size float64) {
	obj.Day.SetFont(file, size)
}

func (obj *GSmallDayCtrl) SetEventText(t string) {
	obj.Event.SetText(t)
}

func (obj *GSmallDayCtrl) SetEventTextColor(c color.Color) {
	obj.Event.SetTextColor(c)
}

func (obj *GSmallDayCtrl) SetEventTextFont(file string, size float64) {
	obj.Event.SetFont(file, size)
}

func (obj *GSmallDayCtrl) SetHolidayText(t string) {
	obj.Holiday.SetText(t)
}

func (obj *GSmallDayCtrl) SetHolidayTextColor(c color.Color) {
	obj.Holiday.SetTextColor(c)
}

func (obj *GSmallDayCtrl) SetHolidayTextFont(file string, size float64) {
	obj.Holiday.SetFont(file, size)
}

func (obj *GSmallDayCtrl) SetLunarText(t string) {
	obj.Lunar.SetText(t)
}

func (obj *GSmallDayCtrl) SetLunarTextColor(c color.Color) {
	obj.Lunar.SetTextColor(c)
}

func (obj *GSmallDayCtrl) SetLunarTextFont(file string, size float64) {
	obj.Lunar.SetFont(file, size)
}