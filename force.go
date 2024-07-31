package lvm2go

type Force bool

func (opt Force) ApplyToVGRemoveOptions(opts *VGRemoveOptions) {
	opts.Force = opt
}

func (opt Force) ApplyToLVRemoveOptions(opts *LVRemoveOptions) {
	opts.Force = opt
}

func (opt Force) ApplyToVGReduceOptions(opts *VGReduceOptions) {
	opts.Force = opt
}

func (opt Force) ApplyToPVRemoveOptions(opts *PVRemoveOptions) {
	opts.Force = opt
}

func (opt Force) ApplyToPVCreateOptions(opts *PVCreateOptions) {
	opts.Force = opt
}

func (opt Force) ApplyToArgs(args Arguments) error {
	if opt {
		args.AddOrReplaceAll([]string{"--force"})
	}
	return nil
}
