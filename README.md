# dbg-go

[![CI build](https://github.com/FedeDP/dbg-go/actions/workflows/ci.yml/badge.svg)](https://github.com/FedeDP/dbg-go/actions/workflows/ci.yml)
[![Latest](https://img.shields.io/github/v/release/FedeDP/dbg-go)](https://github.com/FedeDP/dbg-go/releases/latest)
[![Go Coverage](https://github.com/FedeDP/dbg-go/wiki/coverage.svg)](https://raw.githack.com/wiki/FedeDP/dbg-go/coverage.html)

A go tool to work with falcosecurity [drivers build grid](https://github.com/falcosecurity/test-infra/tree/master/driverkit).  
Long term aim is to completely reimplement dbg Makefile and bash scripts in a much more maintenable and testable language.  

Right now, the tool implements, under the `configs` subcmd:
* configs generation (comprehensive of automatic generation from kernel-crawler output)
* configs cleanup
* configs validation
* configs stats

Moreover, under the `s3` subcmd:
* s3 driver stats
* s3 driver cleanup

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
  completion  Generate the autocompletion script for the specified shell
  configs     Work with local dbg configs
  help        Help about any command
  s3          Work with remote s3 bucket

Flags:
  -a, --architecture string           architecture to run against. (default "x86_64")
      --driver-version strings        driver versions to run against.
      --dry-run                       enable dry-run mode.
  -h, --help                          help for dbg-go
  -l, --log-level string              set log verbosity. (default "INFO")
      --repo-root string              test-infra repository root path. (default "/home/federico/Work/dbg-go")
      --target-distro string          target distro to work against. By default tool will work on any supported distro. Can be a regex.
                                      Supported distros: [almalinux,amazonlinux,amazonlinux2,amazonlinux2022,amazonlinux2023,bottlerocket,centos,debian,fedora,minikube,talos,ubuntu].
      --target-kernelrelease string   target kernel release to work against. By default tool will work on any kernel release. Can be a regex.
      --target-kernelversion string   target kernel version to work against. By default tool will work on any kernel version. Can be a regex.

Use "dbg-go [command] --help" for more information about a command.
```

As you can see, global options basically reimplement all [dbg Makefile filters](https://github.com/falcosecurity/test-infra/blob/master/driverkit/Makefile).

## Build

A simple `make build` in the project root folder is enough.

## Test

Given the project aims at making our dbg code testable, there are already quite a few tests implemented.  
To run them, a simple `make test` issued from project root folder is enough.

## Release artifacts

Using `goreleaser`, multiple artifacts are attached to each github release; among them, you can find executables for arm64 and amd64.

## Examples

<details>
  <summary>Fetch stats about local dbg configs for all supported driver versions by test-infra, for host architecture</summary>
  ```bash
  ./dbg-go configs stats --repo-root test-infra
  ```
</details>

<details>
  <summary>Fetch stats about remote drivers for 5.0.1+driver driver version, for host architecture</summary>
  ```bash
  ./dbg-go s3 stats --driver-version 5.0.1+driver
  ```
</details>

<details>
  <summary>Validate local configs for 5.0.1+driver driver version, for aarch64</summary>
  ```bash
  ./dbg-go configs validate --driver-version 5.0.1+driver --architecture aarch64
  ```
</details>

<details>
  <summary>Generate configs for all supported driver versions by test-infra from kernel-crawler output, for host architecture</summary>
  ```bash
  ./dbg-go configs generate --repo-root test-infra --auto
  ```
</details>