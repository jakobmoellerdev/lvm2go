package lvm2go

type RefreshDevices bool

func (opt RefreshDevices) ApplyToDevUpdateOptions(opts *DevUpdateOptions) {
	opts.RefreshDevices = opt
}

func (opt RefreshDevices) ApplyToArgs(args Arguments) error {
	if opt {
		args.AddOrReplace("--refresh")
	}
	return nil
}
