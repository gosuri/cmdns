package cmdns_test

import (
	"fmt"

	"github.com/gosuri/cmdns"
	"github.com/spf13/cobra"
)

var runFunc = func(cmd *cobra.Command, args []string) { fmt.Println("run", cmd.Name()) }

func Example() {
	ovrclk := &cobra.Command{Use: "ovrclk"}
	apps := &cobra.Command{Use: "apps"}
	apps.AddCommand(&cobra.Command{Use: "info", Run: runFunc})
	ovrclk.AddCommand(apps)

	// Enable namespacing
	cmdns.Enable(ovrclk)

	ovrclk.Execute()
}
