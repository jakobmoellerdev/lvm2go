package lvm2go

type RemoveMissing bool

func (opt RemoveMissing) ApplyToArgs(args Arguments) error {
	if opt {
		args.AddOrReplace("--removemissing")
	}
	return nil
}

func (opt RemoveMissing) ApplyToVGReduceOptions(opts *VGReduceOptions) {
	opts.RemoveMissing = opt
}
