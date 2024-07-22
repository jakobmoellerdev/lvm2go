package lvm2go

type ActivationState string

const (
	Activate     ActivationState = "y"
	Deactivate   ActivationState = "n"
	AutoActivate ActivationState = "ay"
)

func (opt ActivationState) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.ActivationState = opt
}

func (opt ActivationState) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.ActivationState = opt
}

func (opt ActivationState) ApplyToArgs(args Arguments) error {
	if opt == "" {
		return nil
	}
	args.AddOrReplaceAll([]string{"--activate", string(opt)})
	return nil
}
