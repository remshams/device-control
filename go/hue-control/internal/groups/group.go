package groups

type GroupAdapter interface {
	All() ([]Group, error)
}

type Group struct {
	id        string
	name      string
	lights    []string
	connected bool
}

func InitGroup(id string, name string, lights []string) Group {
	return Group{
		id,
		name,
		lights,
		true,
	}
}

func (group Group) GetId() string {
	return group.id
}

func (group Group) GetName() string {
	return group.name
}

func (group Group) GetConnected() bool {
	return group.connected
}

func (group Group) GetLightIds() []string {
	return group.lights
}
