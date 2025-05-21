package worker

import "sync"

// Pool is a pool of goroutines which can run tasks concurrently
type Pool struct {
	size      int
	waitGroup sync.WaitGroup
	tasks     chan Task
}

// NewPool creates a worker pool with the given size
func NewPool(size int) *Pool {
	p := &Pool{
		size:  size,
		tasks: make(chan Task, 0),
	}
	p.spawn()
	return p
}

// spawn spawns n goroutines where n is the size of the pool
func (p *Pool) spawn() {
	for i := 0; i < p.size; i++ {
		go func() {
			for task := range p.tasks {
				task()
				p.waitGroup.Done()
			}
		}()
	}
}

// Submit submits a new task for the workers to run
func (p *Pool) Submit(task Task) {
	p.waitGroup.Add(1)
	p.tasks <- task
}

// WaitAll waits until every task submited has finished
func (p *Pool) WaitAll() {
	p.waitGroup.Wait()
}
