# dbg-go

A go tool to work with falcosecurity [drivers build grid](https://github.com/falcosecurity/test-infra/tree/master/driverkit).  
Long term aim is to completely reimplement dbg Makefile and bash scripts in a much more maintenable and testable language.  

Right now, the tool implements:
* configs generation (comprehensive of automatic generation from kernel-crawler output)
* configs cleanup
* configs validation
* configs stats

This is enough to port [`update-dbg` image](https://github.com/falcosecurity/test-infra/tree/master/images/update-dbg) to make use of this tool instead of the currently used bash scripts.  
First benchmarks showed a tremendous perf improvement: old update-dbg scripts took around 50m on my laptop for a single driverversion. The new tool takes ~10s.  
For more info, see https://github.com/falcosecurity/test-infra/pull/1204#issuecomment-1663822663.  

Tracking issue: https://github.com/falcosecurity/test-infra/issues/1221

## CLI options

Multiple CLI options are available; you can quickly check them out with `./dbg-go --help`, or using `--help` on any sub command.  

```
A command line helper tool used by falcosecurity test-infra dbg.

Usage:
  dbg-go [command]

Available Commands:
  cleanup     Cleanup outdated dbg configs
  completion  Generate the autocompletion script for the specified shell
  generate    Generate new dbg configs
  help        Help about any command
  stats       Fetch stats about configs
  validate    Validate dbg configs

Flags:
  -a, --architecture string           architecture to run against. (default "x86_64")
      --driver-version strings        driver versions to run against.
      --dry-run                       enable dry-run mode.
  -h, --help                          help for dbg-go
  -l, --log-level string              set log verbosity. (default "INFO")
      --repo-root string              test-infra repository root path. (default "/home/federico/Work/dbg-go")
      --target-distro string          target distro to work against. By default tool will work on any supported distro. Can be a regex.
      --target-kernelrelease string   target kernel release to work against. By default tool will work on any kernel release. Can be a regex.
      --target-kernelversion string   target kernel version to work against. By default tool will work on any kernel version. Can be a regex.

Use "dbg-go [command] --help" for more information about a command.
```

As you can see, global options basically reimplement all [dbg Makefile filters](https://github.com/falcosecurity/test-infra/blob/master/driverkit/Makefile).  

## Build

A simple `go build` in the project root folder is enough.

## Test

Given the project aims at making our dbg code testable, there are already quite a few tests implemented.  
To run them, a simple `go test ./...` issued from project root folder is enough.

## Release artifacts

Using `goreleaser`, multiple artifacts are attached to each github release; among them, you can find executables for arm64 and amd64.
