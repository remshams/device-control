package cli

import (
	"keylight-control/control"
	"log"

	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:   "keylight-on",
	Short: "Control keylights",
}

func AddDiscoverCommand(keylightControl *control.KeylightControl) {
	command := &cobra.Command{
		Use:   "discover",
		Short: "Discover keylights",
		Run: func(cmd *cobra.Command, args []string) {
			keylightControl.DiscoverKeylights()
		},
	}
	RootCommand.AddCommand(command)
}

func AddSendCommand(keylightControl *control.KeylightControl) {
	command := &cobra.Command{
		Use:   "sendCommand",
		Short: "Send command to keylight",
		Run: func(cmd *cobra.Command, args []string) {
			on, _ := cmd.Flags().GetBool("on")
			brightness := getInt("brightness", cmd)
			temperature := getInt("temperature", cmd)
			err := keylightControl.SendKeylightCommand(control.KeylightCommand{Id: 0, Command: control.LightCommand{On: &on, Brightness: brightness, Temperature: temperature}})
			if err != nil {
				log.Println(err)
			}
		},
	}
	command.Flags().BoolP("on", "o", false, "Is light on")
	command.Flags().IntP("brightness", "b", 0, "Brightness of light in percent")
	command.Flags().IntP("temperature", "t", 0, "Temperature of the light")
	RootCommand.AddCommand(command)
}

func getInt(name string, cmd *cobra.Command) *int {
	value, err := cmd.Flags().GetInt(name)
	if err == nil && value != 0 {
		return &value
	} else {
		return nil
	}
}
