package game

// Provides a queue of functions to defer processing in the world.
type TaskQueue struct {
	q []func()
}

func newQueue() *TaskQueue {
	q := &TaskQueue{}
	q.q = make([]func(), 0)
	return q
}

// Queues a task.
func (q *TaskQueue) Queue(f func()) {
	q.q = append(q.q, f)
}

// Pops a task.
func (q *TaskQueue) Pop() func() {
	if len(q.q) == 0 {
		return nil
	}
	first := q.q[0]
	q.q = q.q[1:]
	return first
}
