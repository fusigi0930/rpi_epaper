//go:build !rpi

package disp

import (
	"sync"
	"fmt"
	"unsafe"
	"image"

	"ggcal/log"

	"github.com/veandco/go-sdl2/sdl"
)

var gVirtWin *sdl.Window = nil
var gVirtRender *sdl.Renderer = nil
var gWg *sync.WaitGroup = nil
var gQuit bool = false
var gWidth int = 0
var gHeight int = 0

func SDLEventLoop() {
	if gVirtRender == nil {
		log.LogService().Infof("wait for quit loop %v, %v\n", gVirtRender, gVirtWin)
		if gWg != nil {
			gWg.Done()
		}
		return
	}

	for !gQuit {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				gQuit = true
				gWg.Done()
				return
			}
		}

		//gVirtRender.Present()
		sdl.Delay(1000 / 20)
	}
}

func EventLoop() {
	SDLEventLoop()
}

func InitWin(w int, h int) error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.LogService().Errorf("initial SDL2 failed: %v\n", err)
		return err
	}

	gVirtWin, err := sdl.CreateWindow(
		"e-Paper virtual screen",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(w),
		int32(h),
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		log.LogService().Errorf("create SDL2 virtual window failed: %v\n", err)
		return err
	}

	gVirtRender, err = sdl.CreateRenderer(gVirtWin, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.LogService().Errorf("create SDL2 virtual render failed: %v\n", err)
		return err	
	}

	gWg = new(sync.WaitGroup)
	gWg.Add(1)
	gQuit = false

	//go sdl_event_handler()
	gVirtRender.SetDrawColor(255, 255, 255, 255)
	gVirtRender.Clear()

	log.LogService().Infof("initial finished....\n")
	gWidth = w
	gHeight = h
	return nil
}

func CloseWin() {
	gQuit = true
	if gWg != nil {
		gWg.Wait()
	}

	if gVirtRender != nil {
		gVirtRender.Destroy()
		gVirtRender = nil
	}
	if gVirtWin != nil {
		gVirtWin.Destroy()
		gVirtWin = nil
	}
	sdl.Quit()
}

/*
// --- L2 繪圖層模擬 ---
// 這裡可以將來自 L3 的繪圖指令轉換為 SDL2 的繪圖指令

// 1. 畫線 (黑色)
renderer.SetDrawColor(0, 0, 0, 255) // 黑色
renderer.DrawLine(50, 50, 200, 150)

// 2. 畫矩形框 (紅色)
renderer.SetDrawColor(255, 0, 0, 255) // 紅色
rect := sdl.Rect{X: 100, Y: 100, W: 100, H: 50}
renderer.DrawRect(&rect)

// 3. 畫實心矩形 (藍色)
renderer.SetDrawColor(0, 0, 255, 255) // 藍色
filledRect := sdl.Rect{X: 250, Y: 50, W: 80, H: 80}
renderer.FillRect(&filledRect)
*/

func DrawLine(x_start int, y_start int, x_end int, y_end int, c int) {
	if gVirtRender == nil {
		return
	}

	switch (c) {
	case 0: //black
		gVirtRender.SetDrawColor(0, 0, 0, 255)
	case 1: //read
		gVirtRender.SetDrawColor(255, 0, 0, 255)
	case 2: //white
		gVirtRender.SetDrawColor(255, 255, 255, 255)
	}

	gVirtRender.DrawLine(int32(x_start), int32(y_start), int32(x_end), int32(y_end))
	
}

func convColor(img *image.RGBA) []uint8{
	buf := make([]uint8, len(img.Pix))
	for i := 0; i < len(img.Pix); i += 4 {
		// img is RGBA
		// sdl in windows is ABGR
		buf[i] = img.Pix[i+3]
		buf[i+1] = img.Pix[i+2]
		buf[i+2] = img.Pix[i+1]
		buf[i+3] = img.Pix[i+0]
	}
	return buf
}

func UpdateScreen(sc *SurfaceContext) error {
	if sc == nil {
		log.LogService().Errorf("surface context is null\n")
		return fmt.Errorf("surface context is null\n")
	}
	if gVirtRender == nil {
		log.LogService().Errorf("virtual render is null\n")
		return fmt.Errorf("virtual render is null\n")	
	}

	texture, err := gVirtRender.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_STREAMING,
		int32(gWidth),
		int32(gHeight),
	)
	if err != nil {
		log.LogService().Errorf("create texture failed: %v\n", err)
		return err
	}
	defer texture.Destroy()
	colorBuf := convColor(sc.canvas)
	pixPtr := unsafe.Pointer(&colorBuf[0])
	err = texture.Update(nil, pixPtr, sc.canvas.Stride)
	if err != nil {
		log.LogService().Errorf("update texture failed: %v\n", err)
		return err
	}

    gVirtRender.SetDrawColor(0, 0, 0, 255) // 黑色背景
    gVirtRender.Clear()

	gVirtRender.Copy(texture, nil, nil)
	gVirtRender.Present()

	return nil
}