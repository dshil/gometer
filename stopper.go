package gometer

// Stopper is used to stop started entities.
type Stopper interface {
	Stop()
}

var _ Stopper = (*stopperFunc)(nil)

type stopperFunc struct {
	stop func()
}

func (s *stopperFunc) Stop() {
	if s.stop != nil {
		s.stop()
	}
}
