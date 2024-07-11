package lvm2go

type SyncAction string

const (
	SyncActionCheck  SyncAction = "check"
	SyncActionRepair SyncAction = "repair"
)

func (opt SyncAction) ApplyToArgs(args Arguments) error {
	args.AddOrReplaceAll([]string{"--syncaction", string(opt)})
	return nil
}

func (opt SyncAction) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.SyncAction = opt
}
