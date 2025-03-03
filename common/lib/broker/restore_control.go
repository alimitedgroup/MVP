package broker

import "sync"

type RestoreStreamControl struct {
	Wg sync.WaitGroup
}

func NewRestoreStreamControl() *RestoreStreamControl {
	return &RestoreStreamControl{}
}

func (r *RestoreStreamControl) Start() {
	r.Wg.Add(1)
}

func (r *RestoreStreamControl) Finish() {
	r.Wg.Done()
}

func (r *RestoreStreamControl) Wait() {
	r.Wg.Wait()
}
