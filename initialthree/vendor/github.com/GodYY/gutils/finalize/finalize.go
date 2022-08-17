package finalize

import "runtime"

type Finalizable interface {
	Finalizer()
}

func SetFinalizer(o Finalizable) {
	runtime.SetFinalizer(o, finalizer)
}

func finalizer(o Finalizable) {
	o.Finalizer()
}

func UnsetFinalizer(o Finalizable) {
	runtime.SetFinalizer(o, nil)
}
