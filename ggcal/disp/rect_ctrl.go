package disp

import(
	"strconv"
	"image/color"

	"ggcal/log"
)

type GRectControl struct {
	Ui *GObject
	Fill bool
}

func (obj *GRectControl) Init(parent GBase) {
	obj.Ui = &GObject{}
	obj.Ui.Init(parent)
}

func (obj *GRectControl) SetPos(x int, y int) {
	obj.Ui.SetPos(x, y)
}

func (obj *GRectControl) GetPos() (int, int) {
	return obj.Ui.GetPos()
}

func (obj *GRectControl) GetAbsPos() (int, int) {
	return obj.Ui.GetAbsPos()
}

func (obj *GRectControl) SetWidth(w int) {
	obj.Ui.SetWidth(w)
}

func (obj *GRectControl) GetWidth() int {
	return obj.Ui.GetWidth()
}

func (obj *GRectControl) SetHeight(h int) {
	obj.Ui.SetHeight(h)
}

func (obj *GRectControl) GetHeight() int {
	return obj.Ui.GetHeight()
}

func (obj *GRectControl) SetText(text string) {
	obj.Ui.SetText(text)
}

func (obj *GRectControl) GetText() string {
	return obj.Ui.GetText()
}

func (obj *GRectControl) SetTextColor(c color.Color) {
	obj.Ui.SetTextColor(c)
}

func (obj *GRectControl) GetTextColor() color.Color {
	return obj.Ui.GetTextColor()
}

func (obj *GRectControl) SetBorderColor(c color.Color) {
	obj.Ui.SetBorderColor(c)
}

func (obj *GRectControl) GetBorderColor() color.Color {
	return obj.Ui.GetBorderColor()
}

func (obj *GRectControl) SetBorder(b BorderType) {
	obj.Ui.SetBorder(b)
}

func (obj *GRectControl) GetBorder() BorderType {
	return obj.Ui.GetBorder()
}

func (obj *GRectControl) SetParent(parent GBase) {
	obj.Ui.SetParent(parent)
}

func (obj *GRectControl) GetParent() GBase {
	return obj.Ui.GetParent()
}

func (obj *GRectControl) Self() GBase {
	return obj
}

func (obj *GRectControl) AddChild(child GBase) {
	obj.Ui.AddChild(child)
}

func (obj *GRectControl) GetChild() []GBase {
	return obj.Ui.GetChild()
}

func (obj *GRectControl) GetSurfaceContext() *SurfaceContext {
	x, y := obj.GetAbsPos()
	w := obj.GetWidth()
	h := obj.GetHeight()
	sc := obj.Ui.SC.AddChild(x, y, w, h)
	return sc
}

func (obj *GRectControl) SetSurfaceContext(sc *SurfaceContext) {
	obj.Ui.SC = sc
}

func (obj *GRectControl) Draw(scs ...*SurfaceContext) {
	if obj.Ui.SC == nil {
		log.LogService().Errorf("no way....\n")
		return
	}

	//sc := obj.GetSurfaceContext()

	obj.Ui.Draw()
}

func (obj *GRectControl) Close() {
	obj.Ui.Close()
}

func (obj *GRectControl) DumpParam(lead string) {
	obj.Ui.DumpParam(lead+"  ")
}

func (obj *GRectControl) UpdateParam(params map[string] string) {
	obj.Ui.UpdateParam(params)
	fill, ok := params["rectfill"]
	if ok {
		b, _ := strconv.ParseBool(fill)
		obj.Fill = b
	}
}