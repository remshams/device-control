package keylight_model

type ViewState string

const (
	Navigate ViewState = "navigate"
	Insert             = "insert"
	InError            = "inError"
)
