package database

type queue []*Result

func (q queue) Len() int {
	return len(q)
}

func (q queue) Less(i, j int) bool {
	return q[i].Score > q[j].Score
}

func (q queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].Index = i
	q[j].Index = j
}

func (q *queue) Push(x interface{}) {

	var (
		e = x.(*Result)
		n = len(*q)
	)

	e.Index = n
	*q = append(*q, e)

}

func (q *queue) Pop() interface{} {

	var (
		old = *q
		n   = len(old)
		e   = old[n-1]
	)

	old[n-1] = nil
	e.Index = -1

	*q = old[0 : n-1]

	return e

}
