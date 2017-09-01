package flood

// NEXT: The score should be more than a simple count

import (
	"container/heap"
	"fmt"
)

// Search for the best next move.
func (g *Game) SearchMultiLevel(moveSoFar int) Buttons {

	tries := 20000

	// start out our "best" as just beyond optimal.
	bestSoFar := &Item{
		path: make(ButtonPath, 999),
		//path:  make(ButtonPath, g.MaxMoves+1),
		score: 2, // TODO: Try score: 0
	}

	// Initialize the priority queue
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// Seed it with an "no moves" state.
	st := g.At.Copy()
	heap.Push(&pq, &Item{state: &st})

	for {
		// pick best score so far.
		bestItem := heap.Pop(&pq).(*Item)

		// If the path we're about to try is equal or longer
		// than our current best, skip it.
		//  best +1 +2     bestSoFar
		//   17  18 19     19  try it
		//   18  19 20     19  skip
		//   19  20 21     19  skip
		//   20  21 22     19  skip
		if (len(bestItem.path) + 2) > len(bestSoFar.path) {
			continue
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
			fmt.Printf("(%3d) Path %s would score %d (this=%d best=%d)\n", len(pq), newItem.path, newItem.score, len(newItem.path), len(bestSoFar.path))

			if nextScore == g.WinScore {
				// Winning path? Record if it beats our best so far.
				if len(newItem.path) < len(bestSoFar.path) {
					bestSoFar = bestItem
				}
			} else if len(newItem.path) == g.MaxMoves {
				// No more moves? Skip.
				continue
			} else if (bestSoFar.score != g.WinScore) && nextScore > bestSoFar.score {
				// If we don't yet have a winning score, and this is better
				// than what we have, record it as best so far.
				// bestSoFar = newItem
			}
			// Record it for future analysis:
			heap.Push(&pq, newItem)
			//fmt.Printf("PUSHING %v %v %s\n", newItem.path, newItem.score, *newItem.state)

		}

		tries--
		if tries == 0 {
			break
		}
	}

	return bestSoFar.path[0]
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
