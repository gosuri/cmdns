package cmdns

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
)

var helpFunc = func(cmd *cobra.Command, args []string) { cmd.Help() }
var runFunc = func(cmd *cobra.Command, args []string) { fmt.Println("run", cmd.Name()) }

func TestNamespace(t *testing.T) {
	ovrclk := &cobra.Command{Use: "ovrclk"}
	apps := &cobra.Command{Use: "apps"}
	apps.AddCommand(&cobra.Command{Use: "info", Run: runFunc})
	ovrclk.AddCommand(apps)

	// Enable namespacing
	Namespace(ovrclk)
	for _, c := range ovrclk.Commands() {
		if c.Name() == "apps:info" {
			return
		}
	}
	t.Fatalf("expected", "apps:info", "in", ovrclk.Commands())
}
