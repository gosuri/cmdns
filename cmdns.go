package cmdns

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// DefaultNamespacer is the default namespacer for the package
	DefaultNamespacer = New()

	// DefaultNamespaceSeparator is the char that seperates commands
	DefaultNamespaceSeparator = ":"
)

// SetOverrideUsageFunc when set to true will overide the command's usage function with the package's usage function that displays namespaces using the default namespacer
func SetOverrideUsageFunc(v bool) *CobraNamespacer {
	return DefaultNamespacer.SetOverrideUsageFunc(v)
}

// Namespace enables namespacing for the command using the DefaultCmdNS
func Namespace(cmd *cobra.Command) error {
	return DefaultNamespacer.Namespace(cmd)
}

// CobraNamespacer is the struct represting the component that namespaces cobra's sucommands
type CobraNamespacer struct {
	// Namespaces is the collection of cobra namespaces
	Namespaces        []*CobraNamespace
	OverrideUsageFunc bool
}

// New returns a new instance of the CobraNamespacer
func New() *CobraNamespacer {
	return &CobraNamespacer{
		Namespaces:        make([]*CobraNamespace, 0),
		OverrideUsageFunc: true,
	}
}

// SetOverrideUsageFunc when set to true will overide the command's usage
// function with the package's usage function that displays namespaces
func (c *CobraNamespacer) SetOverrideUsageFunc(v bool) *CobraNamespacer {
	c.OverrideUsageFunc = v
	return c
}

// Namespace enables namespacing for the command's subcommands
func (c *CobraNamespacer) Namespace(cmd *cobra.Command) error {
	if cmd == nil {
		return errors.New("cmdns: cmd cannot be nil")
	}
	for _, child := range cmd.Commands() {
		n := NewCobraNamespace()
		n.OverrideUsageFunc = c.OverrideUsageFunc
		c.Namespaces = append(c.Namespaces, n)
		if err := n.Namespace(child); err != nil {
			return err
		}
	}
	return nil
}

// CobraNamespace represents a namespace for a command. This is usually then second level command.
type CobraNamespace struct {
	OverrideUsageFunc bool

	cmd      *cobra.Command
	commands []*cobra.Command
}

// NewCobraNamespace returns a new Namespace
func NewCobraNamespace() *CobraNamespace {
	return &CobraNamespace{
		OverrideUsageFunc: true,
		commands:          make([]*cobra.Command, 0),
	}
}

// AvailableCommands returns the namespaced commands that are available
func (n *CobraNamespace) AvailableCommands() []*cobra.Command {
	return n.commands
}

// Command returns the command for the namespace
func (n *CobraNamespace) Command() *cobra.Command {
	return n.cmd
}

// Namespace enables namespacing for a sub-commmand and its immediated children. It returns an error if the command does not have a parent.
func (n *CobraNamespace) Namespace(cmd *cobra.Command) error {
	if !cmd.HasParent() {
		return errors.New("cmdns: command requires a parent")
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
		nc.Use = cmd.Name() + DefaultNamespaceSeparator + c.Use

		// add this command to the root and hide it so it does not show in available commands list
		c.Parent().Parent().AddCommand(&nc)
		c.Hidden = true
		n.commands = append(n.commands, &nc)
	}
	n.cmd = cmd
	return nil
}

// UsageFunc returns the usage function for the command that renders namespaces
func (n *CobraNamespace) UsageFunc() (f func(*cobra.Command) error) {
	return func(*cobra.Command) error {
		err := tmpl(n.Command().Out(), usageTemplate, n)
		if err != nil {
			fmt.Print(err)
		}
		return err
	}
}

var usageTemplate = `{{$ns := .}}{{with .Command}}Usage: {{if .Runnable}}{{.UseLine}}{{if .HasFlags}} [flags]{{end}}{{end}}{{if gt .Aliases 0}}

Aliases:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}

Examples:
  {{ .Example }}{{end}}{{ if .HasAvailableSubCommands}}

Additional commands, use "{{.Parent.CommandPath}} COMMAND --help" for more information about a command.
{{range $ns.AvailableCommands}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{ if .HasLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimRightSpace}}{{end}}{{ if .HasInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimRightSpace}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}
{{end}}
`
