package diff

func myers(a, b []string) []Operation {
	n := len(a)
	m := len(b)
	maxEdits := n + m

	// Compute the least number of edits needed
	xVals := make([]int, 2*maxEdits+1)
	xVals[maxEdits+1] = 0
	trace := make([][]int, 0)
outer:
	for d := 0; d <= maxEdits; d++ {
		for k := -d; k <= d; k += 2 {
			var x int
			if k == -d || (k != d && xVals[maxEdits+(k-1)] < xVals[maxEdits+(k+1)]) {
				x = xVals[maxEdits+(k+1)]
			} else {
				x = xVals[maxEdits+(k-1)] + 1
			}
			y := x - k

			for x < n && y < m && a[x] == b[y] {
				x++
				y++
			}
			xVals[maxEdits+k] = x

			if x >= n && y >= m {
				trace = append(trace, append([]int{}, xVals...))
				break outer
			}
		}
		trace = append(trace, append([]int{}, xVals...))
	}

	// Backtrack to find the actual operations
	x := n
	y := m
	operations := make([]Operation, 0)

	for d := len(trace) - 1; d >= 0; d-- {
		vals := trace[d]
		k := x - y
		var prevK int

		if k == -d || (k != d && vals[maxEdits+(k-1)] < vals[maxEdits+(k+1)]) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}
		prevX := vals[maxEdits+prevK]
		prevY := prevX - prevK

		for x > prevX && y > prevY {
			x--
			y--
			operations = append([]Operation{eqlOp(x+1, y+1, a[x])}, operations...)
		}

		if d > 0 {
			if x == prevX {
				operations = append([]Operation{insOp(prevY+1, b[prevY])}, operations...)
			} else {
				operations = append([]Operation{delOp(prevX+1, a[prevX])}, operations...)
			}
		}
		x = prevX
		y = prevY
	}
	return operations
}
