package extractflood

import (
	"fmt"
	"image"
	"log"
	"os"

	"github.com/pkg/errors"

	// Uncomment the decoders you want to activate.
	"image/color"
	_ "image/png"
	// _ "image/gif"
	// _ "image/jpeg"
)

func ExtractFile(filename string) error {

	// These are the boundaries of the game field.
	// They were found by trial and error on an iPhone SE.
	// Other phone sizes may result in different dimensions.
	//minX, minY := 14, 278
	//maxX, maxY := 625, 1047
	var minX, minY int = 40, 207
	var maxX, maxY int = minX + 561, minY + 561
	var lenX = maxX - minX
	var lenY = maxY - minY

	reader, err := os.Open(filename)
	if err != nil {
		return errors.Wrapf(err, "ExtractFile:")
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()
	bmy := bounds.Max.Y
	// TODO(tlim): Error if bounds != (0,0)-(640,1136)
	fmt.Println(bounds)

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

	var size string
	var sz int
	fmt.Println("table")
	for irow, row := range runtable {
		if irow == 0 {
			continue
		}
		if row > 1000 {
			fmt.Println(irow, row)
			switch irow {
			case 24:
				size = "LARGE"
				sz = 22
			case 32, 64:
				size = "MEDIUM"
				sz = 17
			case 48:
				size = "SMALL"
				sz = 12
			default:
			}
		}
	}
	fmt.Printf("size=%s\n", size)
	widthX := lenX / sz
	widthY := lenY / sz
	fmt.Printf("widthX=%d widthY=%d\n", widthX, widthY)

	for y := sz - 1; y >= 0; y-- {
		for x := 0; x < sz; x++ {
			// pixel start of the square:
			px := minX + (lenX * x / sz) // more accurate than (lenX/sz)*x + minX
			py := minY + (lenY * y / sz) // more accurate than (lenY/sz)*y + minY
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
			fmt.Printf(" %s", letter(c))
		}
		fmt.Println()
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

var color2letter = map[string]string{}
var lastletter int
var letters = "ABCDEF"

func letter(c color.Color) string {
	r, g, b, _ := c.RGBA()
	u := fmt.Sprintf("%04x%04x%04x", r>>12, g>>12, b>>12)
	v, ok := color2letter[u]
	if ok {
		return v
	}
	newletter := letters[lastletter : lastletter+1]
	lastletter++
	color2letter[u] = newletter
	return newletter

	//r, g, b, _ := c.RGBA()
	//return fmt.Sprintf(" (%d,%d,%d)", r, g, b)
}
