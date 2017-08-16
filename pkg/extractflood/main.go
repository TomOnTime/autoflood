package extractflood

import (
	"fmt"
	"image"
	"log"
	"os"

	"github.com/pkg/errors"

	// Uncomment the decoders you want to activate.
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
	minX, minY := 40, 207
	maxX, maxY := minX+559, minY+561

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
	fmt.Println("table")
	for irow, row := range runtable {
		//		if row != 0 {
		//			fmt.Println(irow, row)
		//		}
		if row == 0 {
			continue
		}
		if row > 10 && row > 1000 {
			fmt.Println(irow, row)
			switch irow {
			case 24:
				size = "LARGE"
			case 32, 64:
				size = "MEDIUM"
			case 48:
				size = "SMALL"
			default:
			}
		}
	}

	fmt.Printf("size: %s\n", size)
	return nil

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		if y < minY || y > maxY {
			continue
		}
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if x < minX || x > maxX {
				continue
			}
			if x%44 != 0 {
				continue
			}
			r, g, b, a := m.At(x, y).RGBA()
			fmt.Printf("pixel %04d:%04d %02d %02d %02d %02x\n", x, y, r>>8, g>>8, b>>8, a)
		}
	}

	return nil
}
