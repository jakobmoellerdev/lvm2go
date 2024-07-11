package lvm2go

type Permission string

const (
	PermissionReadOnly  Permission = "r"
	PermissionReadWrite Permission = "rw"
)

func (opt Permission) ApplyToArgs(args Arguments) error {
	args.AddOrReplaceAll([]string{"--permission", string(opt)})
	return nil
}

func (opt Permission) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.Permission = opt
}
