package lvm2go

type CommonOptions struct {
	Devices
	DevicesFile
	Verbose
}

func (opts CommonOptions) ApplyToArgs(args Arguments) error {
	if err := opts.Devices.ApplyToArgs(args); err != nil {
		return err
	}
	if err := opts.DevicesFile.ApplyToArgs(args); err != nil {
		return err
	}
	if err := opts.Verbose.ApplyToArgs(args); err != nil {
		return err
	}

	return nil
}
