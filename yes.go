package lvm2go

type Yes bool

func (opt Yes) ApplyToArgs(args Arguments) error {
	if opt {
		args.AddOrReplace("--yes")
	}
	return nil
}
