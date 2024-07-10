package lvm2go

import (
	"slices"
	"strings"
)

type Arguments interface {
	AddOrReplaceAll(args []string)
	AddOrReplace(args ...string)
	GetType() ArgsType
	GetRaw() []string
}

type Argument interface {
	ApplyToArgs(args Arguments) error
}

type ArgumentGenerator interface {
	AsArgs() (Arguments, error)
}

type args struct {
	raw []string
	typ ArgsType
}

type ArgsType string

const (
	ArgsTypeLVs      ArgsType = "lvs"
	ArgsTypePVs      ArgsType = "pvs"
	ArgsTypeVGs      ArgsType = "vgs"
	ArgsTypeGeneric  ArgsType = "generic"
	ArgsTypeLVCreate ArgsType = "lvcreate"
	ArgsTypeVGCreate ArgsType = "vgcreate"
	ArgsTypeVGRemove ArgsType = "vgremove"
	ArgsTypeLVRemove ArgsType = "lvremove"
)

func NewArgs(typ ArgsType) Arguments {
	return &args{typ: typ}
}

func (opt *args) AddOrReplaceAll(args []string) {
	for i := range args {
		if fi := slices.Index(opt.raw, args[i]); fi > -1 {
			opt.raw[fi] = args[i]
		} else {
			opt.raw = append(opt.raw, args[i])
		}
	}
}

func (opt *args) AddOrReplace(args ...string) {
	opt.AddOrReplaceAll(args)
}

func (opt *args) GetType() ArgsType {
	return opt.typ
}

func (opt *args) GetRaw() []string {
	return opt.raw
}

func (opt *args) String() string {
	return strings.Join(opt.raw, " ")
}
