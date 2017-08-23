package flood

import "github.com/pkg/errors"

// Flood Performs a floodfill on a State, starting at top,left, replacing
// target with sub.
func (state State) ButtonPress(replace Buttons) (int, error) {

	ly := len(state[0]) - 1
	search := state[0][ly]

	if search == replace {
		return 0, errors.Errorf("Search and replace are the same")
	}
	if state[0][ly] == replace {
		return 0, errors.Errorf("Must start in search-colored area")
	}

	//fmt.Printf("Filling at %d,%d find=%v replace with %v\n", 0, ly, search, replace)
	count := state.fill(0, ly, search, replace, 0)

	return count, nil
}

// fill is the workhorse of ButtonPress, flood filling starting
// at x, y, finding search and replacing it with replace. It returns
// the number of elements replaced.
func (state State) fill(x, y int, search, replace Buttons, count int) int {
	//fmt.Printf("fill(%d, %d, %v, %v)", x, y, search, replace)
	// Illegal location:
	if x < 0 || x >= len(state) || y < 0 || y >= len(state) {
		//fmt.Printf(" BOUNDS\n")
		return 0
	}
	// already has the replacement value
	if state[x][y] == replace {
		//fmt.Printf(" GAIN\n")
		return 0
	}
	// is NOT what we are trying to replace
	if state[x][y] != search {
		//fmt.Printf(" NOT\n")
		return 0
	}
	// IS what we are trying to replace
	state[x][y] = replace
	count += 1
	//fmt.Printf(" REPLACED\n")
	// Try adjact positions
	count += state.fill(x, y-1, search, replace, 0) // above
	count += state.fill(x, y+1, search, replace, 0) // below
	count += state.fill(x-1, y, search, replace, 0) // left
	count += state.fill(x+1, y, search, replace, 0) // right
	return count
}
