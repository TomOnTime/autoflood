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
	var maxX, maxY int = minX + 559, minY + 561

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
			r, g, b, a := m.At(x, y).RGBA()
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
				sz = 24
			case 32, 64:
				size = "MEDIUM"
				sz = 32
			case 48:
				size = "SMALL"
				sz = 48
			default:
			}
		}
	}

	fmt.Printf("size: %s\n", size)

	for y := minY + (sz / 2); y < maxY; y = y + sz {
		fmt.Printf("r %04d", y)
		for x := minX + (sz / 2); x < maxX; x = x + sz {
			//r, g, b, a := m.At(x, y).RGBA()
			//fmt.Printf(" %03d:%03d:%03d:%03x", r>>8, g>>8, b>>8, a>>8)
			c := vote(m.At(x, y), m.At(x+1, y), m.At(x, y+1))
			fmt.Printf(" %s", letter(c))
		}
		fmt.Println()
	}

	return nil
}

func vote(a, b, c color.Color) color.Color {
	if a == b && b == c {
		return a
	}
	if a == b {
		return a
	}
	if a == c {
		return a
	}
	if b == c {
		return b
	}
	fmt.Printf("warning: no majority %v %v %v\n", a, b, c)
	return a
}

func letter(c color.Color) string {
	return "c"
}
