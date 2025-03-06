package broker

import "sync"

type IRestoreStreamControlFactory interface {
	Build() IRestoreStreamControl
}

type RestoreStreamControlFactory struct{}

func NewRestoreStreamControlFactory() *RestoreStreamControlFactory {
	return &RestoreStreamControlFactory{}
}

func (r *RestoreStreamControlFactory) Build() IRestoreStreamControl {
	return NewRestoreStreamControl()
}

type IRestoreStreamControl interface {
	Start()
	Finish()
	Wait()
}

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
