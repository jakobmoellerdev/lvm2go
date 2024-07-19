package lvm2go

import (
	"fmt"
)

type SyncAction string

const (
	SyncActionCheck  SyncAction = "check"
	SyncActionRepair SyncAction = "repair"
)

func (opt SyncAction) ApplyToArgs(args Arguments) error {
	if opt == "" {
		return nil
	}
	args.AddOrReplace(fmt.Sprintf("--syncaction=%s", string(opt)))
	return nil
}

func (opt SyncAction) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.SyncAction = opt
}
