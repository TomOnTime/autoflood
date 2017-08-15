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
	minX, minY := 14, 278
	maxX, maxY := 625, 1047

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

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		if y < minY || y > maxY {
			continue
		}
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if x < minX || x > maxX {
				continue
			}
			r, g, b, a := m.At(x, y).RGBA()
			fmt.Println(x, y, r, g, b, a)
		}
	}

	return nil
}
