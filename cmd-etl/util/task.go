package util

import "sync"

// Worker must be implemented by types that want to use the run pool
type Worker interface {
	Work()
}

// Task provides a pool of goroutines that execute queued Worker
// tasks upon submission
type Task struct {
	work chan Worker
	wg   sync.WaitGroup
}

// New creates a work pool
func New(maxGoroutines int) *Task {
	t := Task{
		// Unbuffered to guarantee that work submitted is actually
		// executed after the call to Run returns
		work: make(chan Worker),
	}

	t.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			for w := range t.work {
				w.Work()
			}
			t.wg.Done()
		}()
	}

	return &t
}

// Do starts work for the pool
func (t *Task) Do(w Worker) {
	t.work <- w
}

// Shutdown halts any ongoing jobs
func (t *Task) Shutdown() {
	close(t.work)
	t.wg.Wait()
}
