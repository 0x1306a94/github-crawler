package crawler

import "math/rand"

type workPool struct {
	pool 		[]*work
	resultQueue chan WorkResult
}

func newWorkPool(poolSize int) *workPool {
	p := &workPool{
		pool: make([]*work, 0),
		resultQueue: make(chan WorkResult, 100),
	}
	for i := 0; i < poolSize; i++ {
		w := NewWork()
		w.id = i + 1
		p.pool = append(p.pool, w)
	}
	return p
}

func (p *workPool) poolSzie() int {
	return len(p.pool)
}

func (p *workPool) addTask(t task)  {
	idx := rand.Intn(p.poolSzie())
	work := p.pool[idx]
	work.jobQueue <- t
}

func (p *workPool) run() *workPool {
	for _, w := range p.pool {
		go w.start()
		go func(wp *workPool, w *work) {
			for result := range w.resultQueue {
				p.resultQueue <- result
			}
		}(p, w)
	}
	return p
}

func (p *workPool) stop() *workPool {
	for _, w := range p.pool {
		go w.stop()
	}
	close(p.resultQueue)
	return p
}