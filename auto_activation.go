package lvm2go

import (
	"fmt"
)

const (
	SetAutoActivate   AutoActivation = "y"
	SetNoAutoActivate AutoActivation = "n"
)

type AutoActivation string

func (opt AutoActivation) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.AutoActivation = opt
}

func (opt AutoActivation) ApplyToArgs(args Arguments) error {
	if opt == "" {
		return nil
	}
	args.AddOrReplace(fmt.Sprintf("--setautoactivation=%s", string(opt)))
	return nil
}
