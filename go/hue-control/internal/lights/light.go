package lights

type LightAdapter interface {
	All() ([]Light, error)
}

type Light struct {
	id   string
	name string
	on   bool
	bri  int
	hue  int
	sat  int
}

func InitLight(id string, name string, on bool, bri int, hue int, sat int) Light {
	return Light{
		id:   id,
		name: name,
		on:   on,
		bri:  bri,
		hue:  hue,
		sat:  sat,
	}
}

func (light Light) GetId() string {
	return light.id
}

func (light Light) GetName() string {
	return light.name
}

func (light Light) GetOn() bool {
	return light.on
}

func (light Light) GetBri() int {
	return light.bri
}

func (light Light) GetHue() int {
	return light.hue
}

func (light Light) GetSat() int {
	return light.sat
}
