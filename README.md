# cmdns

cmdns is a go library for [Cobra](https://github.com/spf13/cobra) for adding namespaces to subcommands. Command namespacing is an alternative way to structure sub commands, similar to `rake db:migrate` and `ovrclk apps:create`.

Cobra a popular library for creating powerful modern CLI applications used by [Kubernetes](http://kubernetes.io/), [Docker](https://github.com/docker/distribution), [Parse](https://github.com/ParsePlatform/parse-cli), and many more widely used Go projects.

## Example

```go
ovrclk := &cobra.Command{Use: "ovrclk"}
apps := &cobra.Command{Use: "apps"}
apps.AddCommand(&cobra.Command{Use: "info", Run: runFunc})
ovrclk.AddCommand(apps)

// Enable namespacing
cmdns.Enable(ovrclk)

ovrclk.Execute()
```

The above example will namespace `info` with `apps`. It lets you run:

```sh
$ ovrclk apps:info
```

And, updates the usage function for `ovrclk:apps -h`:

```sh
Available Commands:
  apps:info

Use "ovrclk [command] --help" for more information about a command.
```

To disable overiding usage function:

```go
cmdns.SetOverrideUsageFunc(false)
```
