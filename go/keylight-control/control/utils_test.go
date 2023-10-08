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

func TestShouldUpdateExistingKeylight(t *testing.T) {
	keylights := prepareKeylights()
	numberOfKeylights := len(keylights)
	newName := "Updated"
	newKeylight := keylights[0]
	newKeylight.Metadata.Name = newName
	updatedKeylights, updatedKeylight := UpdateKeylights(keylights, newKeylight)

	assert.Equal(t, numberOfKeylights, len(updatedKeylights))
	assert.Equal(t, updatedKeylight.Metadata.Name, newName)
}

func TestShouldAddNewKeylight(t *testing.T) {
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
}
