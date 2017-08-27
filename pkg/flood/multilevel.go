package flood

import "fmt"

// Search for the best next move.
func (g *Game) SearchMultiLevel() Buttons {

	count := -1
	best := Buttons(99)
	var max int

	for _, b := range [...]Buttons{0, 1, 2, 3, 4, 5} {
		//count = g.try(Buttons(b))
		count = g.try(b)
		//fmt.Println(g)
		fmt.Printf("Button %v would score %d\n", b, count)
		if count > max {
			max = count
			best = b
		}
	}
	fmt.Printf(" Best: %v %d\n", best, max)
	return best
}
