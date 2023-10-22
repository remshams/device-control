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

func RemoveKeylight(keylights []Keylight, id int) ([]Keylight, *Keylight) {
	for i, _ := range keylights {
		keylight := &keylights[i]
		if keylight.Metadata.Id == id {
			return append(keylights[:i], keylights[i+1:]...), keylight
		}
	}
	return keylights, nil
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
