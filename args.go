package lvm2go

type Arguments interface {
	AppendAll(args []string)
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
	ArgsTypeVGs      ArgsType = "vgs"
	ArgsTypeLVCreate ArgsType = "lvcreate"
	ArgsTypeVGCreate ArgsType = "vgcreate"
	ArgsTypeVGRemove ArgsType = "vgremove"
	ArgsTypeLVRemove ArgsType = "lvremove"
)

func NewArgs(typ ArgsType) Arguments {
	return &args{typ: typ}
}

func (opt *args) AppendAll(args []string) {
	opt.raw = append(opt.raw, args...)
}

func (opt *args) GetType() ArgsType {
	return opt.typ
}

func (opt *args) GetRaw() []string {
	return opt.raw
}
