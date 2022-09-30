package notify

import "sync"

type CancellationChannel struct {
	once     sync.Once
	C        chan struct{}
	Reason   error
	callback func()
}

func NewCancellationChannel() *CancellationChannel {

	return &CancellationChannel{
		C: make(chan struct{}),
	}
}

func NewCancellationChannelWithCallback(callback func()) *CancellationChannel {

	return &CancellationChannel{
		C:        make(chan struct{}),
		callback: callback,
	}
}

func (s *CancellationChannel) WaitForCancellation() {
	<-s.C
}

func (s *CancellationChannel) Cancel() {
	s.once.Do(func() {
		close(s.C)
		if s.callback != nil {
			s.callback()
		}
	})
}

func (s *CancellationChannel) CancelWithReason(reason error) {
	s.Cancel()
	s.Reason = reason
}

func (s *CancellationChannel) IsCancelled() bool {

	isCancelled := true

	select {
	case <-s.C:
	default:
		isCancelled = false
	}

	return isCancelled
}
