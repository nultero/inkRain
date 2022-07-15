package main

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

const spc = ' '

var chars = []rune{
	'｡', '｢', '｣', '､', '･', 'ｦ', 'ｧ', 'ｨ', 'ｩ', 'ｪ', 'ｫ', 'ｬ', 'ｭ', 'ｮ', 'ｯ',
	'ｰ', 'ｱ', 'ｲ', 'ｳ', 'ｴ', 'ｵ', 'ｶ', 'ｷ', 'ｸ', 'ｹ', 'ｺ', 'ｻ', 'ｼ', 'ｽ', 'ｾ', 'ｿ',
	'ﾀ', 'ﾁ', 'ﾂ', 'ﾃ', 'ﾄ', 'ﾅ', 'ﾆ', 'ﾇ', 'ﾈ', 'ﾉ', 'ﾊ', 'ﾋ', 'ﾌ', 'ﾍ', 'ﾎ', 'ﾏ',
	'ﾐ', 'ﾑ', 'ﾒ', 'ﾓ', 'ﾔ', 'ﾕ', 'ﾖ', 'ﾗ', 'ﾘ', 'ﾙ', 'ﾚ', 'ﾛ', 'ﾜ', 'ﾝ', 'ﾞ', 'ﾟ',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
	'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
	'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
}

var grays = []tcell.Style{
	tcell.StyleDefault.Background(tcell.ColorBlack.TrueColor()).Foreground(tcell.ColorLightGray),
	tcell.StyleDefault.Background(tcell.ColorBlack.TrueColor()).Foreground(tcell.ColorDarkGray),
	tcell.StyleDefault.Background(tcell.ColorBlack.TrueColor()).Foreground(tcell.ColorDarkSlateGray),
	tcell.StyleDefault.Background(tcell.ColorBlack.TrueColor()).Foreground(tcell.ColorSlateGray),
}

var bgGreen = tcell.NewRGBColor(0, 19, 9)
var flashGreen = tcell.NewRGBColor(0, 45, 37)

var greens = []tcell.Style{
	tcell.StyleDefault.Background(bgGreen).Foreground(tcell.ColorForestGreen),
	tcell.StyleDefault.Background(bgGreen).Foreground(tcell.ColorGreenYellow),
	tcell.StyleDefault.Background(bgGreen).Foreground(tcell.ColorLightSeaGreen),
	tcell.StyleDefault.Background(bgGreen).Foreground(tcell.ColorLimeGreen),
}

var (
	empt        = []rune{}
	flickerbias = 0
	flickerlast = 0
	pollRate    = 45 * time.Millisecond
	paused      = false
)

const minX, minY = 80, 80

type char_t struct {
	r   rune
	sty tcell.Style
}

func init() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	rand.Seed(time.Now().Unix())
}

func main() {
	scr, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err = scr.Init(); err != nil {
		panic(err)
	}

	scr.HideCursor()
	scr.SetStyle(tcell.StyleDefault)
	scr.Clear()

	chartab := initCharTab()
	coords := getCoords("nult.txt", minX, minY)

	quit := make(chan struct{})
	go func() {
		for {
			ev := scr.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEnter:
					if paused {
						paused = false
					} else {
						paused = true
					}

				case tcell.KeyEscape, tcell.KeyCtrlC, tcell.KeyCtrlQ:
					quit <- struct{}{}
					close(quit)
					return
				case tcell.KeyCtrlL:
					scr.Sync()
				}
			case *tcell.EventResize:
				scr.Sync()
			}
		}
	}()

renderloop:
	for {
		select {
		case <-quit:
			break renderloop

		case <-time.After(pollRate):
			if !paused {
				filltab(chartab)
				drawtab(scr, chartab, coords)
			}
		}

		scr.Show()
	}

	scr.Fini()
}

func drawtab(scr tcell.Screen, ct []char_t, coords coordmap) {
	i := 0
	for x := 0; x < minX; x++ {
		for y := 0; y < minY; y++ {
			c := coord{x: x, y: y}
			if _, ok := coords[c]; ok {
				n := rand.Intn(len(greens))
				sty := greens[n]
				if ct[i].r != spc {
					sty = greens[n].Background(flashGreen)
				}

				r := ct[i].r
				if n < 2 {
					r = '|'
				}

				scr.SetContent(
					x, y, r, empt, sty,
				)

			} else {
				scr.SetContent(
					x, y, ct[i].r, empt, ct[i].sty,
				)
			}

			i++
		}
	}
}

func getChar_t() char_t {

	/*
		I like the way that
		roughly 40% chance of
		empty space in the 'rain'
		looks, so the x here is
		the risk of rain
	*/

	r := ' '
	x := rand.Intn(6)
	if x > 2 {
		r = chars[rand.Intn(len(chars))]
	}
	n := rand.Intn(len(grays))
	c := char_t{
		r:   r,
		sty: grays[n],
	}
	return c
}

func initCharTab() []char_t {
	chartab := make([]char_t, minX*minY)
	for i := 0; i < len(chartab); i++ {
		chartab[i] = getChar_t()
	}
	return chartab
}

func filltab(ct []char_t) []char_t {
	for i := len(ct) - 1; i > minX; i-- {
		ct[i] = ct[i-1]
	}

	for i := minX; i > 0; i-- {
		ct[i] = getChar_t()
	}

	return ct
}

func isDrawable(x, y int) bool {
	if x < minX || y < minY {
		return false
	}

	return true
}
