package lvm2go

type Force bool

func (opt Force) ApplyToVGRemoveOptions(opts *VGRemoveOptions) {
	opts.Force = opt
}

func (opt Force) ApplyToArgs(args Arguments) error {
	if opt {
		args.AppendAll([]string{"--force"})
	}
	return nil
}
