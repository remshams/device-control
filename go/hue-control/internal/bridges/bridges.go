package bridges

func FindBridgeById(bridges []Bridge, id string) *Bridge {
	for _, bridge := range bridges {
		if bridge.id == id {
			return &bridge
		}
	}
	return nil
}
