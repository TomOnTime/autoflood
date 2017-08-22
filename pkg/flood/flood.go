package flood

import (
	"fmt"

	"github.com/pkg/errors"
)

// Flood Performs a floodfill on a State, starting at 0,0, replacing
// target with sub.
func (state State) ButtonPress(replace Buttons) error {

	ly := len(state[0]) - 1
	search := state[0][ly]

	if search == replace {
		return errors.Errorf("Search and replace are the same")
	}
	if state[0][ly] == replace {
		return errors.Errorf("Must start in search-colored area")
	}

	fmt.Printf("Filling at %d,%d find=%v replace with %v\n", 0, ly, search, replace)
	state.fill(0, ly, search, replace)

	return nil
}

func (state State) fill(x, y int, search, replace Buttons) {
	fmt.Printf("fill(%d, %d, %v, %v)", x, y, search, replace)
	if x < 0 || x >= len(state) || y < 0 || y >= len(state) {
		fmt.Printf(" BOUNDS\n")
		return
	}
	if state[x][y] != search {
		fmt.Printf(" NOT\n")
		return
	}
	state[x][y] = replace
	fmt.Printf(" REPLACED\n")
	state.fill(x, y-1, search, replace) // above
	state.fill(x, y+1, search, replace) // below
	state.fill(x-1, y, search, replace) // left
	state.fill(x+1, y, search, replace) // right
}
