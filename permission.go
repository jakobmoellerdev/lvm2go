package lvm2go

import (
	"fmt"
)

type Permission string

const (
	PermissionReadOnly  Permission = "r"
	PermissionReadWrite Permission = "rw"
)

func (opt Permission) ApplyToArgs(args Arguments) error {
	if opt == "" {
		return nil
	}
	args.AddOrReplace(fmt.Sprintf("--permission=%s", string(opt)))
	return nil
}

func (opt Permission) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.Permission = opt
}
