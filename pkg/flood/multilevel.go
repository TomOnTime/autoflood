package flood

import "fmt"

// Search for the best next move.
func (g *Game) SearchMultiLevel() Buttons {

	count := -1
	best := Buttons(99)
	var max int

	for _, b := range [...]Buttons{0, 1, 2, 3, 4, 5} {
		count = g.trymulti(b)
		fmt.Printf("Button %v would score %d\n", b, count)
		if count > max {
			max = count
			best = b
		}
	}
	fmt.Printf(" Best: %v %d\n", best, max)
	return best
}

// try attempts a button press without modifying the game.
func (g *Game) trymulti(b Buttons) int {
	st := g.At.Copy()
	_, err := st.ButtonPress(b)
	if err != nil {
		return 0
	}

	var nb int = 6 // number of buttons

	second := b - 1
	if second < 0 {
		second = Buttons(nb - 1)
	}
	count, err := st.ButtonPress(second)
	if err != nil {
		return 0
	}
	return count
}
