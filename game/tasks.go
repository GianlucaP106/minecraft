package game

// Provides a queue of functions to defer processing in the world.
type Queue[T any] struct {
	q []*T
}

func newQueue[T any]() *Queue[T] {
	q := &Queue[T]{}
	q.q = make([]*T, 0)
	return q
}

// Queues a task.
func (q *Queue[T]) Push(el *T) {
	q.q = append(q.q, el)
}

// Pops a task.
func (q *Queue[T]) Pop() *T {
	if len(q.q) == 0 {
		return nil
	}

	first := q.q[0]
	q.q = q.q[1:]
	return first
}
