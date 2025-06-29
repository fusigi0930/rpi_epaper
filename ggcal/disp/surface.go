package disp

import(
	"os"
	"fmt"
	"errors"
	"image"
	"image/draw"
	"image/color"
	"io/ioutil"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"

	"ggcal/log"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/font/opentype"
)

type TTFFont struct {
	Face font.Face
	Size float64
	Path string
}

type SurfaceContext struct {
	width	int
	height	int
	canvas	*image.RGBA
}

var gFonts map[string]*TTFFont

func NewSurface(w int, h int) *SurfaceContext {
	log.LogService().Infoln("create new surface")
	c := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(c, c.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)

	return &SurfaceContext{width: w, height: h, canvas: c}
}

func (sc *SurfaceContext) AddChild(x int, y int, w int, h int) *SurfaceContext {
	img := sc.canvas.SubImage(image.Rect(x, y, x+w-1, y+h-1))
	child, ok := img.(*image.RGBA)
	if !ok {
		log.LogService().Errorf("create sub image failed\n")
		return nil
	}
	
	return &SurfaceContext{width: w, height: h, canvas: child }
	//return sc
}

func (sc *SurfaceContext) SetPixel(x int, y int, c color.Color) {
	if x < 0 || y < 0 {
		log.LogService().Debugf("out of range x: %d, y: %d\n", x, y)
		return
	}
	
	sc.canvas.Set(x, y, c)
}

func (sc *SurfaceContext) GetPixel(x int, y int) color.Color {
	if x < 0 || x >= sc.width || y < 0 || y >= sc.height {
		return nil
	}

	return sc.canvas.At(x, y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (sc *SurfaceContext) _drawLine (x1 int, y1 int, x2 int, y2 int, c color.Color, interval int) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx := interval
	if x1 >= x2 {
		sx = -interval
	}
	sy := interval
	if y1 >= y2 {
		sy = -interval
	}
	err := dx - dy

	//log.LogService().Debugf("line: %d, %d, %d, %d\n", x1, y1, x2, y2)
	for {
		sc.SetPixel(x1, y1, c)
		if abs(x2-x1) < interval && abs(y2-y1) < interval {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

func (sc *SurfaceContext) DrawLine(x1 int, y1 int, x2 int, y2 int, c color.Color, style ...BorderType) {
	if len(style) != 0 && style[0] == LINE_STYLE_DOTLINE {
		sc._drawLine(x1, y1, x2, y2, c, 3)
	} else {
		sc._drawLine(x1, y1, x2, y2, c, 1)
	}
}

func (sc *SurfaceContext) DrawDotLine(x1 int, y1 int, x2 int, y2 int, c color.Color) {
	sc._drawLine(x1, y1, x2, y2, c, 3)
}

func (sc *SurfaceContext) DrawRect(x int, y int, w int, h int, c color.Color) {
	if w <= 0 || h <= 0 {
		return
	}

	log.LogService().Debugf("rect: %d, %d, %d, %d\n", x, y, x+w-1, y+h-1)
	sc.DrawLine(x, y, x+w-1, y, c)
	sc.DrawLine(x, y, x, y+h-1, c)
	sc.DrawLine(x+w-1, y, x+w-1, y+h-1, c)
	sc.DrawLine(x, y+h-1, x+w-1, y+h-1, c)
}

func (sc *SurfaceContext) DrawFilledRect(x int, y int, w int, h int, c color.Color) {
	if w <= 0 || h <= 0 {
		return
	}

	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			sc.SetPixel(i, j, c)
		}
	}
}

func (sc *SurfaceContext) CleanToColor(c color.Color) {
	draw.Draw(sc.canvas, sc.canvas.Bounds(), image.NewUniform(c), image.Point{}, draw.Src)
}

func (sc *SurfaceContext) DrawImage(x int, y int, file string) {
	_, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) {
		log.LogService().Errorf("the file %s is not exist\n", file)
		return
	}

	f, err := os.Open(file)
	if err != nil {
		log.LogService().Errorf("load file: %s failed\n", file)
		return
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.LogService().Errorf("decode image failed\n")
		return
	}
	
	draw.Draw(sc.canvas, img.Bounds().Add(image.Point{x, y}), img, image.Point{}, draw.Src)
}

func (sc *SurfaceContext) LoadFont(file string, size float64, dpi float64) (*TTFFont, error) {
	//log.LogService().Debugf("loading font...\n")
	if gFonts == nil {
		gFonts = make(map[string]*TTFFont)
	}

	_, err := os.Stat(file)
	if errors.Is(err, os.ErrNotExist) {
		log.LogService().Errorf("the file %s is not exist\n", file)
		return nil, fmt.Errorf("file %s is not exist\n", file)
	}

	key := fmt.Sprintf("%s,%d,%d", file, int(size), int(dpi))
	loaded_font, ok := gFonts[key]
	if ok {
		return loaded_font, nil
	}

	fontByte, err := ioutil.ReadFile(file)
	if err != nil {
		log.LogService().Errorf("read file %s failed\n", file)
		return nil, fmt.Errorf("read file %s failed\n", file)
	}

	procFont, err := opentype.Parse(fontByte)
	if err != nil {
		log.LogService().Errorf("parse ttf failed\n")
		return nil, fmt.Errorf("parse ttf failed\n")
	}

	face, err := opentype.NewFace(procFont, &opentype.FaceOptions{
		Size: size,
		DPI: dpi,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.LogService().Errorf("create ttf face failed\n")
		return nil, fmt.Errorf("create ttf face failed\n")
	}

	newload_font := &TTFFont{Face: face, Size: size, Path: file}
	gFonts[key] = newload_font

	//log.LogService().Debugf("loaded font...\n")
	return newload_font, nil
}

func (sc *SurfaceContext) CloseFont(ttf *TTFFont) {
	if ttf.Face != nil {
		ttf.Face.Close()
	}
}

func (sc *SurfaceContext) DrawText(x int, y int, text string, ttf *TTFFont, c color.Color) {
	if ttf == nil || ttf.Face == nil {
		log.LogService().Errorf("the font pointer is null\n")
		return
	}

	dr := &font.Drawer{
		Dst: sc.canvas,
		Src: image.NewUniform(c),
		Face: ttf.Face,
		Dot: fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y),},
	}

	dr.DrawString(text)
}

func (sc *SurfaceContext) StringPixel(text string, ttf *TTFFont) (int, error) {
	if ttf == nil || ttf.Face == nil {
		log.LogService().Errorf("the font pointer is null\n")
		return 0, fmt.Errorf("font infomrat error\n")
	}

	dr := &font.Drawer{
		Dst: sc.canvas,
		Src: image.NewUniform(color.Black),
		Face: ttf.Face,
		Dot: fixed.Point26_6{X: fixed.I(0), Y: fixed.I(0),},
	}

	textWidth := dr.MeasureString(text)
	intWidth := int(textWidth >> 6)
	return intWidth, nil
}

func (sc *SurfaceContext) SavePNG(file string) error {
	f, err := os.Create(file)
	if err != nil {
		log.LogService().Errorf("create file: %s failed: %v\n", file, err)
		return err
	}
	defer f.Close()
	if err := png.Encode(f, sc.canvas); err != nil {
		log.LogService().Errorf("failed to save to png file: %v\n", err)
		return err		
	}

	log.LogService().Infof("save png file: %s\n", file)
	return nil
}