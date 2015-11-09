# cmdns [![GoDoc](https://godoc.org/github.com/gosuri/cmdns?status.svg)](https://godoc.org/github.com/gosuri/cmdns) [![Build Status](https://travis-ci.org/gosuri/cmdns.svg?branch=master)](https://travis-ci.org/gosuri/cmdns)

cmdns is a go library for [Cobra](https://github.com/spf13/cobra) for name-spacing a command's immediate children. Command name-spacing is an alternative way to structure sub commands, similar to `rake db:migrate` and `ovrclk apps:create`.

Cobra is a popular library for creating powerful modern CLI applications used by [Kubernetes](http://kubernetes.io/), [Docker](https://github.com/docker/distribution), [rkt](https://github.com/coreos/rkt), [Parse](https://github.com/ParsePlatform/parse-cli), and many more widely used Go projects.

## Rationale

Name spacing improves a command line's usability by limiting the depth of its subcommands tree and establishes a clear separation between the command, subcommand and its arguments. Complex commands that have over three levels tend to be harder to remember, especially with arguments. For example, the command to edit a data bag item in chef, a popular infrastructure tool, is `knife data bag edit dogs tibetanspaniel`. There is really no way to tell where the command stops and the arguments begin. This leads to confusion, and generally is hard to repeat. A name-spaced representation could be `knife databag:edit dogs tibetanspaniel`. Its clear here that `knife databag:edit` edits the `tibetanspaniel` item in the data bag `dogs`.

## Usage Example

```go
ovrclk := &cobra.Command{Use: "ovrclk"}
apps := &cobra.Command{Use: "apps"}
apps.AddCommand(&cobra.Command{Use: "info", Run: runFunc})
ovrclk.AddCommand(apps)

// Enable namespacing
cmdns.Namespace(ovrclk)

ovrclk.Execute()
```

The above example will name-space `info` with `apps`. It lets you run:

```sh
$ ovrclk apps:info
```

And, updates the usage function for `ovrclk:apps -h`:

```sh
Available Commands:
  apps:info

Use "ovrclk [command] --help" for more information about a command.
```

To disable overriding usage function:

```go
cmdns.SetOverrideUsageFunc(false)
```
