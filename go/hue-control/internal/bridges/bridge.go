package bridges

type BridgesStore interface {
	Save(bridges []Bridge) error
	Load() ([]Bridge, error)
}

type Bridge struct {
	apiKey string
}

func InitBridge(apiKey string) Bridge {
	return Bridge{apiKey: apiKey}
}
