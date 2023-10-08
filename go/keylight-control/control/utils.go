package control

func UpdateKeylight(keylights []Keylight, newKeylight Keylight) []Keylight {
	oldKeylight := FindKeylightWithId(keylights, newKeylight.Metadata.Id)
	if oldKeylight != nil {
		oldKeylight.Metadata = newKeylight.Metadata
	} else {
		newKeylight.Metadata.Id = len(keylights)
		keylights = append(keylights, newKeylight)
	}
	return keylights
}

func FindKeylightWithId(keylights []Keylight, keylightId int) *Keylight {
	for i := range keylights {
		keylight := &keylights[i]
		if keylight.Metadata.Id == keylightId {
			return keylight
		}
	}
	return nil
}
