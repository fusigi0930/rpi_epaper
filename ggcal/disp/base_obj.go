package disp

import(
	"fmt"
	"strconv"
	"strings"
	"image/color"

	"ggcal/log"
)

type GBase interface {
	Init(GBase)
	SetPos(int, int)
	GetPos() (int, int)
	GetAbsPos() (int, int)
	SetWidth(int)
	GetWidth() int
	SetHeight(int)
	GetHeight() int
	SetText(string)
	GetText() string
	SetTextColor(color.Color)
	GetTextColor() color.Color
	SetBorderColor(color.Color)
	GetBorderColor() color.Color
	SetBorder(BorderType)
	GetBorder() BorderType
	Draw(...*SurfaceContext)
	GetSurfaceContext() *SurfaceContext
	SetSurfaceContext(*SurfaceContext)
	Close()
	DumpParam(string)
	UpdateParam(map[string] string)

	SetParent(GBase)
	GetParent() GBase
	AddChild(GBase)
	GetChild() []GBase

}

var gMapUniId map[string]GBase

type BorderType int
const (
	BORDER_NONE BorderType	= 0
	BORDER_ALL			= BORDER_LEFT | BORDER_RIGHT | BORDER_TOP | BORDER_BOTTOM
	BORDER_LEFT			= 0x01
	BORDER_RIGHT		= 0x02
	BORDER_TOP			= 0x04
	BORDER_BOTTOM		= 0x08
	BORDER_VERT 		= BORDER_LEFT | BORDER_RIGHT
	BORDER_HORI			= BORDER_TOP | BORDER_BOTTOM
	LINE_STYLE_LINE		= 0x00
	LINE_STYLE_DOTLINE	= 0x10
)

type GObject struct {
	X int
	Y int
	W int
	H int
	abs_x int
	abs_y int
	Text string
	Tcolor color.Color
	Border BorderType
	Bcolor color.Color
	SC *SurfaceContext
	
	Parent GBase
	Child []GBase
}

func (obj *GObject) Init(parent GBase) {
	obj.Parent = parent
	obj.X = 0
	obj.Y = 0
	obj.Tcolor = color.Black
	obj.Bcolor = color.Black
	//obj.SC = obj.Parent.GetSurfaceContext()

	if parent == nil {
		obj.abs_x = obj.X
		obj.abs_y = obj.Y
	} else {
		obj.abs_x, obj.abs_y = parent.GetAbsPos()
		obj.abs_x += obj.X
		obj.abs_y += obj.Y
	}

}

func (obj *GObject) SetPos(x int, y int) {
	obj.X = x
	obj.Y = y

	if obj.Parent == nil {
		obj.abs_x = x
		obj.abs_y = y
	} else {
		obj.abs_x, obj.abs_y = obj.GetParent().GetAbsPos()
		obj.abs_x += x
		obj.abs_y += y
	}
}

func (obj *GObject) GetPos() (int, int) {
	return obj.X, obj.Y
}

func (obj *GObject) GetAbsPos() (int, int) {
	return obj.abs_x, obj.abs_y
}

func (obj *GObject) SetWidth(w int) {
	obj.W = w
}

func (obj *GObject) GetWidth() int {
	return obj.W
}

func (obj *GObject) SetHeight(h int) {
	obj.H = h
}

func (obj *GObject) GetHeight() int {
	return obj.H
}

func (obj *GObject) SetText(text string) {
	obj.Text = text
}

func (obj *GObject) GetText() string {
	return obj.Text
}

func (obj *GObject) SetTextColor(c color.Color) {
	obj.Tcolor = c
}

func (obj *GObject) GetTextColor() color.Color {
	return obj.Tcolor
}

func (obj *GObject) SetBorderColor(c color.Color) {
	obj.Bcolor = c
}

func (obj *GObject) GetBorderColor() color.Color {
	return obj.Bcolor
}

func (obj *GObject) SetBorder(b BorderType) {
	obj.Border = b
}

func (obj *GObject) GetBorder() BorderType {
	return obj.Border
}

func (obj *GObject) SetParent(parent GBase) {
	obj.Parent = parent

	if obj.Parent == nil {
		obj.abs_x = obj.X
		obj.abs_y = obj.Y

	} else {
		obj.abs_x, obj.abs_y = obj.GetParent().GetAbsPos()
		obj.abs_x += obj.X
		obj.abs_y += obj.Y
	}
}

func (obj *GObject) GetParent() GBase {
	return obj.Parent
}

func (obj *GObject) Self() GBase {
	return obj
}

func (obj *GObject) AddChild(child GBase) {
	if child == nil {
		return
	}

	obj.Child = append(obj.Child, child)
}

func (obj *GObject) GetChild() []GBase {
	return obj.Child
}

func (obj *GObject) GetSurfaceContext() *SurfaceContext {
	x, y := obj.GetAbsPos()
	w := obj.GetWidth()
	h := obj.GetHeight()
	sc := obj.SC.AddChild(x, y, w, h)
	return sc
}

func (obj *GObject) SetSurfaceContext(sc *SurfaceContext) {
	obj.SC = sc
}

func (obj *GObject) Draw(scs ...*SurfaceContext) {
	if obj.Parent == nil && obj.SC == nil {
		log.LogService().Errorf("= =||||\n")
		return
	} else if obj.Parent == nil {
		obj.SC.CleanToColor(color.White)
	} else {
		var extsc *SurfaceContext = nil
		if len(scs) > 0 {
			extsc = scs[0]
		}
		sc := obj.GetSurfaceContext()
		if extsc != nil {
			sc = extsc
		}
		border := obj.GetBorder()
		c := obj.GetBorderColor()
		x, y := obj.GetAbsPos()
		w := obj.GetWidth()
		h := obj.GetHeight()

		if border & BORDER_ALL == BORDER_ALL {
			sc.DrawRect(x, y, w-1, h-1, c)
		} else {
			if border & BORDER_LEFT == BORDER_LEFT {
				sc.DrawLine(x, y, x, y+h-2, c, border & LINE_STYLE_DOTLINE)
			}
			if border & BORDER_RIGHT == BORDER_RIGHT {
				sc.DrawLine(x+w-2, y, x+w-2, y+h-2, c, border & LINE_STYLE_DOTLINE)
			}
			if border & BORDER_TOP == BORDER_TOP {
				sc.DrawLine(x, y, x+w-2, y, c, border & LINE_STYLE_DOTLINE)
			}
			if border & BORDER_BOTTOM == BORDER_BOTTOM {
				sc.DrawLine(x, y+h-2, x+w-2, y+h-2, c, border & LINE_STYLE_DOTLINE)
			}
		}
	}

	for _, child := range obj.GetChild() {
		child.Draw(scs...)
	}

	if obj.Parent == nil {
		UpdateScreen(obj.SC)
	}
}

func (obj *GObject) Close() {
	for _, child := range obj.GetChild() {
		child.Close()
	}
}

func (obj *GObject) DumpParam(lead string) {
	log.LogService().Debugf("%sx: %d\n", lead, obj.X)
	log.LogService().Debugf("%sy: %d\n", lead, obj.Y)
	log.LogService().Debugf("%sabs_x: %d\n", lead, obj.abs_x)
	log.LogService().Debugf("%sabs_y: %d\n", lead, obj.abs_y)
	log.LogService().Debugf("%swidth: %d\n", lead, obj.W)
	log.LogService().Debugf("%sheight: %d\n", lead, obj.H)
	log.LogService().Debugf("%sborder: %d\n", lead, obj.Border)
	log.LogService().Debugf("%sborder_color: %v\n", lead, obj.Bcolor)
	log.LogService().Debugf("%stext_color: %v\n", lead, obj.Tcolor)
	log.LogService().Debugf("%stext: %s\n", lead, obj.Text)
	b := obj.GetSurfaceContext().canvas.Bounds()
	log.LogService().Debugf("%sbounds: %d, %d, %d, %d\n", lead, b.Min.X, b.Min.Y, b.Max.X-1, b.Max.Y-1)

	for _, c := range obj.Child {
		c.DumpParam(lead+"  ")
		log.LogService().Debugf("\n")
	}
}

func strToColor(s string) color.Color {
	var c color.Color
	if s[:1] == "#" {
		v, _ := strconv.ParseUint(s[1:7], 16, 32)
		c = color.RGBA{A:255, R:uint8(v >> 16), G:uint8((v >> 8) & 0xff), B:uint8(v & 0xff)}
	} else {
		switch(s) {
		case "black":
			c = color.Black
		case "white":
			c = color.White
		case "red":
			c = color.RGBA{R:255, A:255}
		case "green":
			c = color.RGBA{G:255, A:255}
		case "blue":
			c = color.RGBA{B:255, A:255}
		case "yellow":
			c = color.RGBA{R:255, G:255, A:255}
		case "gray":
			c = color.RGBA{R:127, G:127, B:127, A:255}
		
		}
	}
	return c
}

func (obj *GObject) UpdateParam(params map[string] string) {
	x, _ := strconv.Atoi(params["x"])
	y, _ := strconv.Atoi(params["y"])
	w, _ := strconv.Atoi(params["width"])
	h, _ := strconv.Atoi(params["height"])
	obj.SetPos(x, y)
	obj.SetWidth(w)
	obj.SetHeight(h)
	t, ok := params["text"]
	if ok {
		obj.SetText(t)
	}
	
	bo, ok := params["border"]
	if ok {
		types := strings.Split(bo, "|")
		var b BorderType = BORDER_NONE
		for _, t := range types {
			switch(t) {
			case "none":
				b = BORDER_NONE
			case "all":
				b |= BORDER_ALL
			case "left":
				b |= BORDER_LEFT
			case "right":
				b |= BORDER_RIGHT
			case "top":
				b |= BORDER_TOP
			case "bottom":
				b |= BORDER_BOTTOM
			case "vert":
				b |= BORDER_VERT
			case "hori":
				b |= BORDER_HORI
			case "line":
				b |= LINE_STYLE_LINE
			case "dotline":
				b |= LINE_STYLE_DOTLINE
			}
		}
		obj.SetBorder(b)
	}

	bc, ok := params["border_color"]
	if ok {
		obj.SetBorderColor(strToColor(bc))
	}

	tc, ok := params["text_color"]
	if ok {
		obj.SetTextColor(strToColor(tc))
	}
}

func SetUniqueId(id string, obj GBase) error {
	if gMapUniId == nil {
		gMapUniId = make(map[string]GBase)
	}
	if _, ok := gMapUniId[id]; !ok {
		gMapUniId[id]=obj
	} else {
		return fmt.Errorf("the key had already exist.\n", id)
	}
	return nil
}

func GetControlById(id string) GBase {
	if _, ok := gMapUniId[id]; ok {
		return gMapUniId[id]
	} 
	return nil
}