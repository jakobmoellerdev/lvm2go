package lvm2go

type CommonOptions struct {
	Devices
	DevicesFile
	Verbose
	RequestConfirm
}

func (opts CommonOptions) ApplyToArgs(args Arguments) error {
	for _, arg := range []Argument{
		opts.Devices,
		opts.DevicesFile,
		opts.Verbose,
		opts.RequestConfirm,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}

type RequestConfirm bool

func (opt RequestConfirm) ApplyToArgs(args Arguments) error {
	if !opt {
		args.AddOrReplace("--yes")
	}
	return nil
}

type Verbose bool

func (opt Verbose) ApplyToArgs(args Arguments) error {
	if opt {
		args.AddOrReplace("--verbose")
	}
	return nil
}
