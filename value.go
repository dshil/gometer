package gometer

import "sync/atomic"

type value struct {
	val int64
}

func (v *value) Inc() {
	v.Add(1)
}

func (v *value) Add(val int64) {
	atomic.AddInt64(&v.val, val)
}

func (v *value) Value() int64 {
	return atomic.LoadInt64(&v.val)
}

func (v *value) Set(val int64) {
	atomic.StoreInt64(&v.val, val)
}
