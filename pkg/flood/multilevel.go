package flood

import (
	"container/heap"
	"fmt"
)

// Search for the best next move.
func (g *Game) SearchMultiLevel() Buttons {

	// Initialize the priority queue
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// Seed it with an "no moves" state.
	st := g.At.Copy()
	heap.Push(&pq, &Item{state: &st})

	for {
		// pick best score so far.
		bestItem := heap.Pop(&pq).(*Item)

		// if max depth: return this {path,score}
		if len(bestItem.path) == g.MaxMoves {
			return bestItem.path[0]
			// BUG???
		}

		// generate a {past + button, score} for each button (remove any loss)
		for _, b := range [...]Buttons{0, 1, 2, 3, 4, 5} {

			// Score it.
			nextState := (*bestItem.state).Copy() // Make a copy
			nextScore := nextState.multitry(b)    // Press the button
			if nextScore == 0 {
				// 0 means "this path loses (was more than MaxMoves moves)
				continue
			}
			//fmt.Printf("nextScore=%d g.WinScore=%d\n", nextScore, g.WinScore)

			// Store it for later analysis.
			newItem := &Item{
				path:  bestItem.path.Plus(b),
				score: nextScore,
				state: &nextState,
			}
			fmt.Printf("(%d) Path %s would score %d\n", len(pq), newItem.path, newItem.score)
			if nextScore == g.WinScore {
				fmt.Printf("RETURN first of %s\n", newItem.path)
				return newItem.path[0]
			}
			//fmt.Printf("PUSHING %v %v %s\n", newItem.path, newItem.score, *newItem.state)
			heap.Push(&pq, newItem)

		}
	}

}

// try attempts a button press without modifying the game.
func (st *State) multitry(b Buttons) Score {
	_, err := st.ButtonPress(b)
	if err != nil {
		return 0
	}

	var nb int = 6 // number of buttons

	// Now that we've press "b", we can count how big the
	// area is by pressing any button other than "b" on a copy.
	// Pick any button other than b
	second := b - 1
	if second < 0 {
		second = Buttons(nb - 1)
	}
	cp := st.Copy()
	count, err := cp.ButtonPress(second)
	if err != nil {
		return 0
	}

	return count
}
