package lvm2go

import (
	"fmt"
)

type ActivationMode string

const (
	ActivationModePartial  ActivationMode = "partial"
	ActivationModeDegraded ActivationMode = "degraded"
	ActivationModeComplete ActivationMode = "complete"
)

func (opt ActivationMode) ApplyToArgs(args Arguments) error {
	if opt == "" {
		return nil
	}
	args.AddOrReplace(fmt.Sprintf("--activationmode=%s", string(opt)))
	return nil
}

func (opt ActivationMode) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.ActivationMode = opt
}
