package util

type Rollbacks struct {
	canceled          bool
	rollbackFunctions []func()
}

func (r *Rollbacks) Add(f func()) {
	r.rollbackFunctions = append(r.rollbackFunctions, f)
}

func (r *Rollbacks) Do() {
	if r.canceled {
		return
	}

	for i := len(r.rollbackFunctions) - 1; i >= 0; i-- {
		r.rollbackFunctions[i]()
	}
}

func (r *Rollbacks) Cancel() {
	r.canceled = true
}
