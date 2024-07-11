package lvm2go

type WipeSignatures bool

func (opt WipeSignatures) ApplyToArgs(args Arguments) error {
	if opt {
		args.AddOrReplace("--wipesignatures")
	}
	return nil
}
