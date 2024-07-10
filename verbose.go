package lvm2go

type Verbose bool

func (opt Verbose) ApplyToArgs(args Arguments) error {
	if opt {
		args.AddOrReplace("--verbose")
	}
	return nil
}
