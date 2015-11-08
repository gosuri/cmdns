package cmdns

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

// DefaultCmdNS is the default CmdNS for the package
var DefaultCmdNS = New()

// SetOverrideUsageFunc when set to true will overide the command's usage function with the package's usage function that displays namespaces using the DefaultCmdNS
func SetOverrideUsageFunc(v bool) *CmdNS {
	return DefaultCmdNS.SetOverrideUsageFunc(v)
}

// Enable enables namespacing for the command using the DefaultCmdNS
func Enable(cmd *cobra.Command) error {
	return DefaultCmdNS.Enable(cmd)
}

// CmdNS is the struct represting the component that namespaces sucommands
type CmdNS struct {
	// Namespaces
	Namespaces        []*Namespace
	OverrideUsageFunc bool
}

// New returns a new instance of the CmdNS
func New() *CmdNS {
	return &CmdNS{
		Namespaces:        make([]*Namespace, 0),
		OverrideUsageFunc: true,
	}
}

// SetOverrideUsageFunc when set to true will overide the command's usage
// function with the package's usage function that displays namespaces
func (c *CmdNS) SetOverrideUsageFunc(v bool) *CmdNS {
	c.OverrideUsageFunc = v
	return c
}

// Enable enables namespacing for the command's subcommands
func (c *CmdNS) Enable(cmd *cobra.Command) error {
	if cmd == nil {
		return errors.New("cmdns: cmd cannot be nil")
	}
	for _, child := range cmd.Commands() {
		n := NewNamespace()
		n.OverrideUsageFunc = c.OverrideUsageFunc
		c.Namespaces = append(c.Namespaces, n)
		if err := n.Enable(child); err != nil {
			return err
		}
	}
	return nil
}

// Namespace represents a namespace for a command
type Namespace struct {
	OverrideUsageFunc bool

	cmd      *cobra.Command
	commands []*cobra.Command
}

// NewNamespace returns a new Namespace
func NewNamespace() *Namespace {
	return &Namespace{
		OverrideUsageFunc: true,
		commands:          make([]*cobra.Command, 0),
	}
}

// AvailableCommands returns the namespaced commands that are available
func (n *Namespace) AvailableCommands() []*cobra.Command {
	return n.commands
}

// Command returns the command for the namespace
func (n *Namespace) Command() *cobra.Command {
	return n.cmd
}

// Enable enables namespacing for a sub-commmand and its immediated children.
// It returns an error if the command does not have a parent
func (n *Namespace) Enable(cmd *cobra.Command) error {
	if !cmd.HasParent() {
		return errors.New("cmdns: command is required a parent")
	}

	// Do not bind if there are not available sub commands
	if !cmd.HasAvailableSubCommands() {
		return nil
	}

	if n.OverrideUsageFunc {
		cmd.SetUsageFunc(n.UsageFunc())
	}

	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() {
			continue
		}
		// copy the command add it to the root command with a prefix of its parent.
		nc := *c
		nc.Use = cmd.Name() + ":" + c.Use
		c.Parent().AddCommand(&nc)

		// hide the command so it does not show in available commands list
		nc.Hidden = true
		n.commands = append(n.commands, &nc)
	}
	n.cmd = cmd
	return nil
}

// UsageFunc returns the usage function for the command that renders namespaces
func (n *Namespace) UsageFunc() (f func(*cobra.Command) error) {
	return func(*cobra.Command) error {
		err := tmpl(n.Command().Out(), usageTemplate, n)
		if err != nil {
			fmt.Print(err)
		}
		return err
	}
}

var usageTemplate = `{{$ns := .}}{{with .Command}}Usage:{{if .Runnable}}
  {{.UseLine}}{{if .HasFlags}} [flags]{{end}}{{end}}{{if gt .Aliases 0}}

Aliases:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}

Examples:
{{ .Example }}{{end}}{{ if .HasAvailableSubCommands}}

Available Commands:{{range $ns.AvailableCommands}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{ if .HasLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimRightSpace}}{{end}}{{ if .HasInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimRightSpace}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasSubCommands }}

Use "{{.Parent.CommandPath}} [command] --help" for more information about a command.{{end}}{{end}}
`
