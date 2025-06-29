//go:build rpi

package disp

/*
#include "DEV_Config.h"
#include "EPD_12in48b.h"
*/
import "C"

import (
	"sync"
	"fmt"
	"image"
	"time"

	"ggcal/log"
)

var gWg *sync.WaitGroup = nil
var gQuit bool = false
var gWidth int = 0
var gHeight int = 0

func EventLoop(drawSignal <-chan struct{}) {
	for !gQuit {
		select {
		case <-drawSignal:
			GetRootScreen().Draw()
		default:
		}
		//gVirtRender.Present()
		time.Sleep(100 * time.Millisecond)
	}
}

func InitWin(w int, h int) error {
	log.LogService().Infof("init epaper module\n")
	C.DEV_ModuleInit()
	if ret := C.EPD_12in48B_Init(); ret != 0 {
		log.LogService().Error("init epaper failed\n")
		return fmt.Errorf("init epaper error code: %d\n", ret)
	}

	C.EPD_12in48B_Clear()
	time.Sleep(3 * time.Second)
	
	log.LogService().Infof("initial finished....\n")
	return nil
}

func CloseWin() {
	gQuit = true
	if gWg != nil {
		gWg.Wait()
	}

}

func DrawLine(x_start int, y_start int, x_end int, y_end int, c int) {
	
}

func convColor(img *image.RGBA) []uint8{
	buf := make([]uint8, len(img.Pix))
	for i := 0; i < len(img.Pix); i += 4 {
		// img is RGBA
		buf[i] = img.Pix[i+3]
		buf[i+1] = img.Pix[i+2]
		buf[i+2] = img.Pix[i+1]
		buf[i+3] = img.Pix[i+0]
	}
	return buf
}

func splitImageToRedBlackArray(img *image.RGBA) ([]uint8, []uint8) {
	l := len(img.Pix)/4 / 8
	if (len(img.Pix)/4) % 8 != 0 {
		l++
	}

	black := make([]uint8, l)
	red := make([]uint8, l)
	for i := range black {
		black[i] = 0xff
		red[i] = 0xff 
	}

	for i := 0; i < len(img.Pix); i += 4 {
		index := i / 4
		addr := index / 8
		if img.Pix[i] <= 0x40 && img.Pix[i+1] <= 0x40 && img.Pix[i+2] <= 0x40 {
			d := black[addr]
			var mask uint8 = 0x80 >> (index % 8)
			mask = mask ^ 0xff
			black[addr] = d & mask
		} else if img.Pix[i] >= 0xd0 && img.Pix[i+1] <= 0x10 && img.Pix[i+2] <= 0x10 {
			d := red[addr]
			var mask uint8 = 0x80 >> (index % 8)
			mask = mask ^ 0xff
			red[addr] = d & mask
		}
	}

	return black, red
}

func UpdateScreen(sc *SurfaceContext) error {
	log.LogService().Infof("translating color image to epaper format...\n")
	b, r := splitImageToRedBlackArray(sc.canvas)
	bPtr := C.CBytes(b)
	rPtr := C.CBytes(r)
	log.LogService().Infof("updating epaper screen...\n")
	C.EPD_12in48B_Display((*C.UBYTE)(bPtr), (*C.UBYTE)(rPtr))
	return nil
}