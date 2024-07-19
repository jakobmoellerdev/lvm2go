package lvm2go

import (
	"fmt"
)

type Type string

const (
	TypeLinear     Type = "linear"
	TypeStriped    Type = "striped"
	TypeMirrored   Type = "mirrored"
	TypeRAID0      Type = "raid0"
	TypeRAID1      Type = "raid1"
	TypeRAID4      Type = "raid4"
	TypeRAID5      Type = "raid5"
	TypeRAID6      Type = "raid6"
	TypeRAID10     Type = "raid10"
	TypeThin       Type = "thin"
	TypeCache      Type = "cache"
	TypeWriteCache Type = "writecache"
	TypePool       Type = "cache-pool"
	TypeThinPool   Type = "thin-pool"
	TypeVDO        Type = "vdo"
	TypeVDOPool    Type = "vdo-pool"
)

func (opt Type) ApplyToArgs(args Arguments) error {
	if opt == "" {
		return nil
	}
	args.AddOrReplace(fmt.Sprintf("--type=%s", string(opt)))
	return nil
}

func (opt Type) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.Type = opt
}
