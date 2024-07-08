package lvm2go

import (
	"strings"
)

var (
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
	DefaultVGsColumnOptions = ColumnOptions{
		"vg_uuid",
		"vg_name",
		"vg_size",
		"vg_free",
		"pv_count",
		"lv_count",
		"snap_count",
		"vg_attr",
	}
	DefaultPVsColumnOptions = ColumnOptions{
		"pv_uuid",
		"pv_name",
		"pv_fmt",
		"pv_size",
		"pv_free",
		"pv_attr",
		"pv_tags",
		"vg_name",
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
	args.AppendAll([]string{"--options", optionsString})
	return nil
}
