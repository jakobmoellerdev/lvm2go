package lvm2go

type AutoActivationFromReport string

const (
	AutoActivationFromReportEnabled  AutoActivationFromReport = "enabled"
	AutoActivationFromReportDisabled AutoActivationFromReport = ""
)

func (opt AutoActivationFromReport) True() bool {
	return opt == AutoActivationFromReportEnabled
}
