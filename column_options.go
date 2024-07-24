package lvm2go

import (
	"strings"
)

var (
	DefaultLVsColumnOptions = ColumnOptions{
		"lv_all",
	}
	DefaultVGsColumnOptions = ColumnOptions{
		"vg_all",
	}
	DefaultPVsColumnOptions = ColumnOptions{
		"pv_all",
	}
)

type ColumnOptions []string

func (opt ColumnOptions) ApplyToLVsOptions(opts *LVsOptions) {
	opts.ColumnOptions = opt
}

func (opt ColumnOptions) ApplyToVGsOptions(opts *VGsOptions) {
	opts.ColumnOptions = opt
}

func (opt ColumnOptions) ApplyToArgs(args Arguments) error {
	var optionsString string
	if len(opt) > 0 {
		optionsString = strings.Join(opt, ",")
	} else {
		switch args.GetType() {
		case ArgsTypeVGs:
			optionsString = strings.Join(DefaultVGsColumnOptions, ",")
		case ArgsTypeLVs:
			optionsString = strings.Join(DefaultLVsColumnOptions, ",")
		case ArgsTypePVs:
			optionsString = strings.Join(DefaultPVsColumnOptions, ",")
		}
	}
	args.AddOrReplaceAll([]string{"--options", optionsString})
	return nil
}
