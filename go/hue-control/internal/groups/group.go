package groups

type GroupAdapter interface {
	All() ([]Group, error)
	Set(group Group) error
}

type Group struct {
	adapter   GroupAdapter
	id        string
	name      string
	lights    []string
	connected bool
	on        bool
}

func InitGroup(adapter GroupAdapter, id string, name string, lights []string, on bool) Group {
	return Group{
		adapter,
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

func (group *Group) SetOn(on bool) {
	group.on = on
}

func (group Group) SendGroup() error {
	return group.adapter.Set(group)
}
