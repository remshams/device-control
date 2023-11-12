package groups

type GroupAdapter interface {
	All() ([]Group, error)
}

type Group struct {
	id        string
	name      string
	lights    []string
	connected bool
	on        bool
}

func InitGroup(id string, name string, lights []string, on bool) Group {
	return Group{
		id,
		name,
		lights,
		true,
		on,
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

func (group Group) GetOn() bool {
	return group.on
}
