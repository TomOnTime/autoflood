package flood

import (
	"bytes"
	"fmt"
	"image"
	"os"

	"github.com/pkg/errors"

	// Uncomment the decoders you want to activate.
	"image/color"
	_ "image/png"
	// _ "image/gif"
	// _ "image/jpeg"
)

type Buttons uint

const (
	Purple Buttons = iota
	Blue
	Green
	Yellow
	Red
	Pink
)

var letters = "ABCDEF"

func (b Buttons) String() string {
	//return letters[lastletter : lastletter+1]
	//return fmt.Sprintf("%d", b)
	return letters[b : b+1]
}

type Game struct {
	Image      image.Image
	Level      string
	Size       int
	At         [][]Buttons
	minX, minY int
	maxX, maxY int
	lenX, lenY int
}

func (g *Game) LoadImage(filename string) (err error) {
	reader, err := os.Open(filename)
	if err != nil {
		return errors.Wrapf(err, "ExtractFile:")
	}
	defer reader.Close()

	g.Image, _, err = image.Decode(reader)

	return
}

func (g *Game) IdentifyLevel() (err error) {

	// These are the boundaries of the game field.
	// They were found by trial and error on an iPhone SE.
	// Other phone sizes may result in different dimensions.
	//minX, minY := 14, 278
	//maxX, maxY := 625, 1047
	var minX, minY int = 40, 207
	var maxX, maxY int = minX + 561, minY + 561
	var lenX = maxX - minX
	var lenY = maxY - minY

	m := g.Image

	bounds := m.Bounds()
	bmy := bounds.Max.Y
	// TODO(tlim): Error if bounds != (0,0)-(640,1136)
	fmt.Println("bounds = ", bounds)

	// Find a run of a color, record the run length.
	//runtable := map[int]int{}
	var runtable [999]int

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		if y < minY || y > maxY {
			continue
		}
		first := true
		var ar, ag, ab, aa uint32
		var run int
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if x < minX || x > maxX {
				continue
			}
			r, g, b, a := m.At(x, bmy-y).RGBA()
			//fmt.Printf("pixel %04x:%04d %02d %02d %02d %02x\n", x, y, r, g, b, a)
			if first {
				first = false
				run = 1
				ar, ag, ab, aa = r, g, b, a
			} else {
				if r == ar && g == ag && b == ab && a == aa {
					run++
				} else {
					runtable[run&(0xffffff-3)]++ // Round to the nearest multiple of 4
					// fmt.Println(x, y, r, g, b, a, run)
					run = 1
					ar, ag, ab, aa = r, g, b, a
				}
			}
		}
	}

	// style rows   pixels    total
	// small:   12   46.75    561
	// medium   17   33       561
	// large    22   25.5     561

	fmt.Println("table")
	for irow, row := range runtable {
		if irow == 0 {
			continue
		}
		if row > 1000 {
			fmt.Println(irow, row)
			switch irow {
			case 24:
				g.Level = "LARGE"
				g.Size = 22
			case 32, 64:
				g.Level = "MEDIUM"
				g.Size = 17
			case 48:
				g.Level = "SMALL"
				g.Size = 12
			default:
			}
		}
	}
	fmt.Printf("boardsize=%s\n", g.Level)

	g.minX, g.minY = minX, minY
	g.maxX, g.maxY = maxX, maxY
	g.lenX, g.lenY = lenX, lenY
	return
}

func (g *Game) ExtractGrid() (err error) {
	sz := g.Size
	m := g.Image
	bmy := m.Bounds().Max.Y

	g.At = make([][]Buttons, sz)
	for i := range g.At {
		g.At[i] = make([]Buttons, sz)
	}

	widthX := g.lenX / g.Size
	widthY := g.lenY / g.Size
	fmt.Printf("widthX=%d widthY=%d\n", widthX, widthY)

	// populate Grid
	for y := sz - 1; y >= 0; y-- {
		for x := 0; x < sz; x++ {
			// pixel start of the square:
			px := g.minX + (g.lenX * x / sz) // more accurate than (lenX/sz)*x + minX
			py := g.minY + (g.lenY * y / sz) // more accurate than (lenY/sz)*y + minY
			// Mid-point
			mx := px + (widthX / 2)
			my := py + (widthY / 2)
			cl := []color.Color{
				m.At(mx-1, bmy-(my-1)), m.At(mx, bmy-(my-1)), m.At(mx+1, bmy-(my-1)),
				m.At(mx-1, bmy-(my)), m.At(mx, bmy-(my)), m.At(mx+1, bmy-(my)),
				m.At(mx-1, bmy-(my+1)), m.At(mx, bmy-(my+1)), m.At(mx+1, bmy-(my+1)),
			}
			c, err := vote(cl)
			//fmt.Printf("%2d:%-2d %3d:%3d %3d:%3d %s [[%v]]\n", x, y, px, py, mx, my, letter(c), c)
			//fmt.Printf("                         %v\n", cl)
			if err != nil {
				fmt.Printf("\nERROR: %s (c=%v)\n", err, c)
			}
			let := letter(c)
			//fmt.Printf(" %s", let)
			g.At[x][y] = let
		}
		//fmt.Println()
	}

	return nil
}

func vote(cl []color.Color) (color.Color, error) {
	var err error
	tally := map[string]int{}

	var max int
	var maxc color.Color

	for _, c := range cl {
		r, g, b, _ := c.RGBA()
		u := fmt.Sprintf("%04x%04x%04x", r>>12, g>>12, b>>12)
		tally[u]++
		if tally[u] > max {
			max = tally[u]
			maxc = c
		}
	}

	cm := 0
	for _, v := range tally {
		if v == max {
			cm++
		}
	}

	if cm != 1 {
		err = fmt.Errorf("no majority in %v ==== %v\n", cl, tally)
	}

	return maxc, err
}

var color2letter = map[string]Buttons{}
var lastletter Buttons

func letter(c color.Color) Buttons {
	r, g, b, _ := c.RGBA()
	u := fmt.Sprintf("%04x%04x%04x", r>>12, g>>12, b>>12)
	v, ok := color2letter[u]
	if ok {
		return v
	}
	color2letter[u] = lastletter
	ret := lastletter
	lastletter++
	return ret

	//newletter := letters[lastletter : lastletter+1]
	//lastletter++
	//color2letter[u] = newletter
	//return newletter

	//r, g, b, _ := c.RGBA()
	//return fmt.Sprintf(" (%d,%d,%d)", r, g, b)
}

func (g *Game) String() string {
	b := bytes.NewBufferString("")
	for y := g.Size - 1; y >= 0; y-- {
		for x := 0; x < g.Size; x++ {
			b.WriteString(fmt.Sprintf(" %v", g.At[x][y]))
		}
		b.WriteString("\n")
	}
	return b.String()
}
