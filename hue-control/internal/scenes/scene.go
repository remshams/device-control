package scenes

type Scene struct {
	id        string
	name      string
	groupId   string
	sceneType string
}

type SceneAdapter interface {
	All(groupId string) ([]Scene, error)
}

func InitScene(id string, name string, groupId string, sceneType string) Scene {
	return Scene{
		id:        id,
		name:      name,
		groupId:   groupId,
		sceneType: sceneType,
	}
}

func (s Scene) Name() string {
	return s.name
}

func (s Scene) GroupId() string {
	return s.groupId
}

func (s Scene) Id() string {
	return s.id
}

func (s Scene) SceneType() string {
	return s.sceneType
}
