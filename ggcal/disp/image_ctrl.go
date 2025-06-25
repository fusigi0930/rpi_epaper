package disp

import(
	"image/color"
	"os"
	"errors"

	"ggcal/log"
)

type GImageControl struct {
	Ui *GObject
	File string
}

func (obj *GImageControl) Init(parent GBase) {
	obj.Ui = &GObject{}
	obj.Ui.Init(parent)
}

func (obj *GImageControl) SetPos(x int, y int) {
	obj.Ui.SetPos(x, y)
}

func (obj *GImageControl) GetPos() (int, int) {
	return obj.Ui.GetPos()
}

func (obj *GImageControl) GetAbsPos() (int, int) {
	return obj.Ui.GetAbsPos()
}

func (obj *GImageControl) SetWidth(w int) {
	obj.Ui.SetWidth(w)
}

func (obj *GImageControl) GetWidth() int {
	return obj.Ui.GetWidth()
}

func (obj *GImageControl) SetHeight(h int) {
	obj.Ui.SetHeight(h)
}

func (obj *GImageControl) GetHeight() int {
	return obj.Ui.GetHeight()
}

func (obj *GImageControl) SetText(text string) {
	obj.Ui.SetText(text)
}

func (obj *GImageControl) GetText() string {
	return obj.Ui.GetText()
}

func (obj *GImageControl) SetTextColor(c color.Color) {
	obj.Ui.SetTextColor(c)
}

func (obj *GImageControl) GetTextColor() color.Color {
	return obj.Ui.GetTextColor()
}

func (obj *GImageControl) SetBorderColor(c color.Color) {
	obj.Ui.SetBorderColor(c)
}

func (obj *GImageControl) GetBorderColor() color.Color {
	return obj.Ui.GetBorderColor()
}

func (obj *GImageControl) SetBorder(b BorderType) {
	obj.Ui.SetBorder(b)
}

func (obj *GImageControl) GetBorder() BorderType {
	return obj.Ui.GetBorder()
}

func (obj *GImageControl) SetParent(parent GBase) {
	obj.Ui.SetParent(parent)
}

func (obj *GImageControl) GetParent() GBase {
	return obj.Ui.GetParent()
}

func (obj *GImageControl) Self() GBase {
	return obj
}

func (obj *GImageControl) AddChild(child GBase) {
	obj.Ui.AddChild(child)
}

func (obj *GImageControl) GetChild() []GBase {
	return obj.Ui.GetChild()
}

func (obj *GImageControl) GetSurfaceContext() *SurfaceContext {
	x, y := obj.GetAbsPos()
	w := obj.GetWidth()
	h := obj.GetHeight()
	sc := obj.Ui.SC.AddChild(x, y, w, h)
	return sc
}

func (obj *GImageControl) SetSurfaceContext(sc *SurfaceContext) {
	obj.Ui.SC = sc
}

func (obj *GImageControl) Draw(scs ...*SurfaceContext) {
	if obj.Ui.SC == nil {
		log.LogService().Errorf("no way....\n")
		return
	}

	sc := obj.GetSurfaceContext()
	x, y := obj.GetAbsPos()
	sc.DrawImage(x, y, obj.File)
	obj.Ui.Draw()
}

func (obj *GImageControl) SetImage(file string) {
	_, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) {
		log.LogService().Errorf("the file %s is not exist\n", file)
		return
	}

	obj.File = file
}

func (obj *GImageControl) Close() {
	obj.Ui.Close()
}

func (obj *GImageControl) DumpParam(lead string) {
	log.LogService().Debugf("%simage: %s\n", lead, obj.File)
	obj.Ui.DumpParam(lead+"  ")
}

func (obj *GImageControl) UpdateParam(params map[string] string) {
	obj.Ui.UpdateParam(params)
	img, ok := params["imgfile"]
	if ok {
		obj.SetImage(img)
	}

	uniid, ok := params["id"]
	if ok {
		SetUniqueId(uniid, obj)
	}
}