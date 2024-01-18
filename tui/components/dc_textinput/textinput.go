package dc_textinput

import (
	"fmt"
	"github.com/remshams/device-control/tui/styles"

	"github.com/charmbracelet/bubbles/textinput"
)

func CreateTextInputModel() textinput.Model {
	model := textinput.New()
	model.TextStyle = styles.TextAccentColor
	return model
}

func CreateTextInputView(model textinput.Model, label string, unit string) string {
	return fmt.Sprintf("%s %s%s", label, model.View(), styles.TextAccentColor.Render(unit))
}
