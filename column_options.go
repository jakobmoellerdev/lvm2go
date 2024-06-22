package lvm2go

import (
	"strings"
)

var (
	DefaultVGsColumnOptions = ColumnOptions{
		"vg_uuid",
		"vg_name",
		"vg_size",
		"vg_free",
	}
	DefaultLVsColumnOptions = ColumnOptions{
		"lv_uuid",
		"lv_name",
		"lv_full_name",
		"lv_path",
		"lv_size",
		"lv_kernel_major",
		"lv_kernel_minor",
		"origin",
		"origin_size",
		"pool_lv",
		"lv_tags",
		"lv_attr",
		"vg_name",
		"data_percent",
		"metadata_percent",
		"pool_lv",
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
		}
	}
	args.AppendAll([]string{"--options", optionsString})
	return nil
}
