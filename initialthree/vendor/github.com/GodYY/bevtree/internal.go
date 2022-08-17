package bevtree

type internal interface {
	__()
}

type internalImpl struct{}

func (internalImpl) __() {}
