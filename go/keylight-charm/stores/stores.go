package stores

type Layout struct {
	Width  int
	Height int
}

var LayoutStore = Layout{
	Width:  0,
	Height: 0,
}
