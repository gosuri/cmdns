package main

import (
	"fmt"

	"github.com/gosuri/cmdns"
	"github.com/spf13/cobra"
)

var helpFunc = func(cmd *cobra.Command, args []string) { cmd.Help() }
var runFunc = func(cmd *cobra.Command, args []string) { fmt.Println("run", cmd.Name()) }

func main() {
	hugo := &cobra.Command{
		Use:   "hugo",
		Short: "Hugo is a very fast static site generator",
		Long:  "A Fast and Flexible Static Site Generator built with love by spf13 and friends in Go",
	}

	newCmd := &cobra.Command{
		Use:   "new",
		Short: "Create new content for your site",
		Long:  "Create a new content file and automatically set the date and title. It will guess which kind of file to create based on the path provided.",
		Run:   runFunc,
	}
	hugo.AddCommand(newCmd)

	newSiteCmd := &cobra.Command{
		Use:   "site",
		Short: "Create a new site in the provided directory",
		Long:  "The new site will have the correct structure, but no content or theme yet",
		Run:   runFunc,
	}
	newCmd.AddCommand(newSiteCmd)

	newThemeCmd := &cobra.Command{
		Use:   "theme",
		Short: "Create a new theme",
		Long:  "Create a new theme (skeleton) called [name] in the current directory. New theme is a skeleton.",
		Run:   runFunc,
	}
	newCmd.AddCommand(newThemeCmd)

	// listCmd := &cobra.Command{Use: "list", Run: runFunc}
	// listCmd.AddCommand(&cobra.Command{Use: "list", Run: runFunc})
	// listCmd.AddCommand(&cobra.Command{Use: "add", Run: runFunc})
	//hugo.AddCommand(listCmd)

	cmdns.Enable(hugo)
	if err := hugo.Execute(); err != nil {
		panic(err)
	}
}
