package control

func UpdateKeylights(keylights []Keylight, newKeylight Keylight) ([]Keylight, Keylight) {
	updatedKeylight := FindKeylightWithId(keylights, newKeylight.Metadata.Id)
	if updatedKeylight != nil {
		updatedKeylight.Metadata = newKeylight.Metadata
	} else {
		updatedKeylight = &newKeylight
		updatedKeylight.Metadata.Id = len(keylights)
		keylights = append(keylights, *updatedKeylight)
	}
	return keylights, *updatedKeylight
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
