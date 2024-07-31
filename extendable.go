package lvm2go

const (
	ExtendableTrue  Extendable = "extendable"
	ExtendableFalse Extendable = ""
)

type Extendable string

func (e Extendable) True() bool {
	return e == ExtendableTrue
}
