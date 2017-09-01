package flood

type State [][]Buttons

//const EmptyState = State

// Copy returns a deep copy of st.
func (st State) Copy() (dst State) {
	for x, xv := range st {
		dst = append(dst, make([]Buttons, len(xv)))
		for y, yv := range xv {
			dst[x][y] = yv
		}
	}
	return dst
}
