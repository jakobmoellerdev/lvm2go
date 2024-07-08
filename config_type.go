package lvm2go

type ConfigType string

const (
	ConfigTypeFull ConfigType = "full"
)

func (c ConfigType) String() string {
	return string(c)
}

func (c ConfigType) AsArgs() []string {
	return []string{"--typeconfig", c.String()}
}

func (c ConfigType) ApplyToConfigOptions(opts *ConfigOptions) {
	opts.ConfigType = c
}

func (c ConfigType) ApplyToArgs(arguments Arguments) error {
	if len(c) == 0 {
		return nil
	}
	arguments.AppendAll(c.AsArgs())
	return nil
}
