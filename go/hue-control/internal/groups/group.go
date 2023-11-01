package groups

type Group struct {
	id     string
	name   string
	lights []string
}

type GroupAdapter interface {
	All() ([]Group, error)
}
