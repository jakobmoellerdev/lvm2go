package lvm2go

import "fmt"

type RemoveMissing VolumeGroupName

func (opt RemoveMissing) ApplyToArgs(args Arguments) error {
	if opt == "" {
		return nil
	}
	args.AddOrReplace(fmt.Sprintf("--removemissing=%s", opt))
	return nil
}

func (opt RemoveMissing) ApplyToVGReduceOptions(opts *VGReduceOptions) {
	opts.RemoveMissing = opt
}
