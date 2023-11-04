package groups

type GroupAdapter interface {
	All() ([]Group, error)
}

type Group struct {
	id     string
	name   string
	lights []string
}

func (group Group) GetId() string {
	return group.id
}

func (group Group) GetName() string {
	return group.name
}
