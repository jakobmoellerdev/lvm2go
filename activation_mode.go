package lvm2go

type ActivationMode string

const (
	ActivationModePartial  ActivationMode = "partial"
	ActivationModeDegraded ActivationMode = "degraded"
	ActivationModeComplete ActivationMode = "complete"
)

func (opt ActivationMode) ApplyToArgs(args Arguments) error {
	args.AddOrReplaceAll([]string{"--activationmode", string(opt)})
	return nil
}

func (opt ActivationMode) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.ActivationMode = opt
}
