package disp

import(
	"image/color"
	"os"
	"errors"
	"strconv"
	"strings"
	"fmt"

	"ggcal/log"
)

type GTextControl struct {
	Ui *GObject
	FontPath string
	Font *TTFFont
	Size float64
	Inverse bool
}

func (obj *GTextControl) Init(parent GBase) {
	obj.Ui = &GObject{}
	obj.Ui.Init(parent)
	obj.Inverse = false
}

func (obj *GTextControl) SetPos(x int, y int) {
	obj.Ui.SetPos(x, y)
}

func (obj *GTextControl) GetPos() (int, int) {
	return obj.Ui.GetPos()
}

func (obj *GTextControl) GetAbsPos() (int, int) {
	return obj.Ui.GetAbsPos()
}

func (obj *GTextControl) SetWidth(w int) {
	obj.Ui.SetWidth(w)
}

func (obj *GTextControl) GetWidth() int {
	return obj.Ui.GetWidth()
}

func (obj *GTextControl) SetHeight(h int) {
	obj.Ui.SetHeight(h)
}

func (obj *GTextControl) GetHeight() int {
	return obj.Ui.GetHeight()
}

func (obj *GTextControl) SetText(text string) {
	obj.Ui.SetText(text)
}

func (obj *GTextControl) GetText() string {
	return obj.Ui.GetText()
}

func (obj *GTextControl) SetTextColor(c color.Color) {
	obj.Ui.SetTextColor(c)
}

func (obj *GTextControl) GetTextColor() color.Color {
	return obj.Ui.GetTextColor()
}

func (obj *GTextControl) SetBorderColor(c color.Color) {
	obj.Ui.SetBorderColor(c)
}

func (obj *GTextControl) GetBorderColor() color.Color {
	return obj.Ui.GetBorderColor()
}

func (obj *GTextControl) SetBorder(b BorderType) {
	obj.Ui.SetBorder(b)
}

func (obj *GTextControl) GetBorder() BorderType {
	return obj.Ui.GetBorder()
}

func (obj *GTextControl) SetParent(parent GBase) {
	obj.Ui.SetParent(parent)
}

func (obj *GTextControl) GetParent() GBase {
	return obj.Ui.GetParent()
}

func (obj *GTextControl) Self() GBase {
	return obj
}

func (obj *GTextControl) AddChild(child GBase) {
	obj.Ui.AddChild(child)
}

func (obj *GTextControl) GetChild() []GBase {
	return obj.Ui.GetChild()
}

func (obj *GTextControl) GetSurfaceContext() *SurfaceContext {
	x, y := obj.GetAbsPos()
	w := obj.GetWidth()
	h := obj.GetHeight()
	sc := obj.Ui.SC.AddChild(x, y, w, h)
	return sc
}

func (obj *GTextControl) SetSurfaceContext(sc *SurfaceContext) {
	obj.Ui.SC = sc
}

func (obj *GTextControl) Draw(scs ...*SurfaceContext) {
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

	if obj.Inverse {
		pan_color = color.White
		bk_color = obj.GetTextColor()
	}
	
	sc.CleanToColor(bk_color)
	for i, t := range lines {
		sc.DrawText(x+1, y+int(obj.Font.Size+1)*(i+1), t, obj.Font, pan_color)
	}

	obj.Ui.Draw()
}

func (obj *GTextControl) SetFont(file string, size float64) {
	_, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) {
		log.LogService().Errorf("the file %s is not exist\n", file)
		return
	}

	obj.Font, _ = obj.Ui.SC.LoadFont(file, size, 72.0)
}

func (obj *GTextControl) Close() {
	if obj.Font != nil {
		obj.Ui.SC.CloseFont(obj.Font)
		obj.Font = nil
	}

	obj.Ui.Close()
}

func (obj *GTextControl) DumpParam(lead string) {
	log.LogService().Debugf("%sfont: %v\n", lead, obj.Font)
	obj.Ui.DumpParam(lead+"  ")
}

func (obj *GTextControl) UpdateParam(params map[string] string) {
	obj.Ui.UpdateParam(params)
	fs, ok := params["fontsize"]
	var size float64
	if ok {
		size, _ = strconv.ParseFloat(fs, 64)
	}

	ft, ok := params["font"]
	if ok {
		path := GlobalSetting("fontpath")
		f := ft
		if len(path) != 0 {
			f = fmt.Sprintf("%s/%s", path, ft)
		}
		obj.SetFont(f, size)
	}

	uniid, ok := params["id"]
	if ok {
		SetUniqueId(uniid, obj)
	}
}

func (obj *GTextControl) StringPixel() int {
	pxl, _ := obj.GetSurfaceContext().StringPixel(obj.GetText(), obj.Font)
	return pxl
}