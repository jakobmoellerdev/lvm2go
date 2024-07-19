package lvm2go

type Discards string

const (
	DiscardsPassdown   Discards = "passdown"
	DiscardsNoPassdown Discards = "nopassdown"
	DiscardsIgnore     Discards = "ignore"
)

func (opt Discards) ApplyToArgs(args Arguments) error {
	if opt == "" {
		return nil
	}
	args.AddOrReplaceAll([]string{"--discards", string(opt)})
	return nil
}

func (opt Discards) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.Discards = opt
}
