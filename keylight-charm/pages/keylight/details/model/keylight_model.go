package keylight_model

type ViewState string

const (
	Navigate ViewState = "navigate"
	Insert             = "insert"
	isError            = "inError"
)

type CommandStatus string

const (
	NoCommand CommandStatus = "noCommand"
	Success                 = "success"
	Error                   = "error"
)
