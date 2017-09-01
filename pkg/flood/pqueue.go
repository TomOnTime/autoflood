package flood

import (
	"bytes"
	"math"
)

const Win = math.MaxInt16

type ButtonPath []Buttons

type Score uint16

func (bp ButtonPath) Plus(b Buttons) ButtonPath {
	var n ButtonPath
	for _, j := range bp {
		n = append(n, j)
	}
	n = append(n, b)
	return n
}

func (bp ButtonPath) String() string {
	buf := bytes.NewBufferString("")
	for _, b := range bp {
		_, err := buf.WriteString(b.String())
		if err != nil {
			panic(err)
		}
	}
	return buf.String()
}

// An Item is something we manage in a priority queue.
type Item struct {
	path  ButtonPath // Path taken to get here
	score Score      // The score (the priority)
	state *State
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, score so we use greater than here.
	return pq[i].score > pq[j].score
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
//func (pq *PriorityQueue) update(item *Item, value string, score int) {
//	item.value = value
//	item.score = score
//	heap.Fix(pq, item.index)
//}
