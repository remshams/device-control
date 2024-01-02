package control

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepareKeylights() []Keylight {
	return []Keylight{{Metadata: KeylightMetadata{
		Id:   0,
		Name: "First",
		Ip:   net.ParseIP("192.168.2.1"),
		Port: 9999,
	}}}
}

func TestUpdateKeylight(t *testing.T) {
	t.Run("should update existing keylight", func(t *testing.T) {
		keylights := prepareKeylights()
		numberOfKeylights := len(keylights)
		newName := "Updated"
		newKeylight := keylights[0]
		newKeylight.Metadata.Name = newName
		updatedKeylights, updatedKeylight := UpdateKeylights(keylights, newKeylight)

		assert.Equal(t, numberOfKeylights, len(updatedKeylights))
		assert.Equal(t, updatedKeylight.Metadata.Name, newName)
	})

	t.Run("should add new keylight", func(t *testing.T) {
		keylights := prepareKeylights()
		numberOfKeylights := len(keylights)
		newKeylight := Keylight{Metadata: KeylightMetadata{
			Id:   -1,
			Name: "New keylight",
			Ip:   net.ParseIP("192.168.1.1"),
			Port: 9998,
		}}
		updatedKeylights, updatedKeylight := UpdateKeylights(keylights, newKeylight)

		assert.Equal(t, numberOfKeylights+1, len(updatedKeylights))
		assert.Equal(t, 1, updatedKeylight.Metadata.Id)
	})
}

func TestRemoveKeylight(t *testing.T) {
	t.Run("should remove existing keylight", func(t *testing.T) {
		keylights := prepareKeylights()
		keylightToRemove := keylights[0]
		numberOfKeylights := len(keylights)
		keylights, removedKeylight := RemoveKeylight(keylights, keylightToRemove.Metadata.Id)

		assert.Equal(t, numberOfKeylights-1, len(keylights))
		assert.Equal(t, &keylightToRemove, removedKeylight)
	})

	t.Run("should do nothing if keylight does not exist", func(t *testing.T) {
		keylights := prepareKeylights()
		numberOfKeylights := len(keylights)
		keylights, removedKeylight := RemoveKeylight(keylights, 99)

		assert.Equal(t, numberOfKeylights, len(keylights))
		assert.Nil(t, removedKeylight)
	})
}
