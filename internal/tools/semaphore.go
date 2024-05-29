package tools

type Semaphore struct {
	ch    chan struct{}
	limit int
}

func NewSemaphore(limit int) *Semaphore {
	return &Semaphore{
		limit: limit,
		ch:    make(chan struct{}, limit),
	}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.ch
}

func (s *Semaphore) TryAcquire() bool {
	if len(s.ch) >= s.limit {
		return false
	}

	s.Acquire()
	return true
}
