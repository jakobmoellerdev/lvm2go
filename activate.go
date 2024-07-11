package lvm2go

type Activate string

const (
	Y  Activate = "y"
	N  Activate = "n"
	AY Activate = "ay"
)

func (opt Activate) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.Activate = opt
}

func (opt Activate) ApplyToArgs(args Arguments) error {
	if opt == "" {
		return nil
	}
	args.AddOrReplaceAll([]string{"--activate", string(opt)})
	return nil
}
