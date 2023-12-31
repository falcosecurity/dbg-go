# dbg-go

 [![Falco Infra Repository](https://github.com/falcosecurity/evolution/blob/main/repos/badges/falco-infra-blue.svg)](https://github.com/falcosecurity/evolution/blob/main/REPOSITORIES.md#infra-scope) [![Incubating](https://img.shields.io/badge/status-incubating-orange?style=for-the-badge)](https://github.com/falcosecurity/evolution/blob/main/REPOSITORIES.md#incubating) [![License](https://img.shields.io/github/license/falcosecurity/dbg-go?style=for-the-badge)](./LICENSE)

[![CI build](https://github.com/falcosecurity/dbg-go/actions/workflows/ci.yml/badge.svg)](https://github.com/falcosecurity/dbg-go/actions/workflows/ci.yml)
[![Latest](https://img.shields.io/github/v/release/falcosecurity/dbg-go)](https://github.com/falcosecurity/dbg-go/releases/latest)
![Architectures](https://img.shields.io/badge/ARCHS-x86__64%7Caarch64-blueviolet)
[![Go Report Card](https://goreportcard.com/badge/github.com/falcosecurity/dbg-go)](https://goreportcard.com/report/github.com/falcosecurity/dbg-go)

A go tool to work with falcosecurity [drivers build grid](https://github.com/falcosecurity/test-infra/tree/master/driverkit).  

The tool implements, under the `configs` subcmd:
* configs generation (comprehensive of automatic generation from kernel-crawler output)
* configs cleanup
* configs validation
* configs stats
* configs build (using driverkit libraries)

Moreover, under the `drivers` subcmd:
* remote driver stats
* remote driver cleanup
* remote driver publish

## CLI options

Multiple CLI options are available; you can quickly check them out with `./dbg-go --help`, or using `--help` on any sub command.  

```
A command line helper tool used by falcosecurity test-infra dbg.

Usage:
  dbg-go [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  configs     Work with local dbg configs
  drivers     Work with remote drivers bucket
  help        Help about any command

Flags:
  -a, --architecture string           architecture to run against. Supported: [amd64,arm64] (default "amd64")
      --driver-name string            driver name to be used (default "falco")
      --driver-version strings        driver versions to run against.
      --dry-run                       enable dry-run mode.
  -h, --help                          help for dbg-go
  -l, --log-level string              set log verbosity. (default "INFO")
      --repo-root string              test-infra repository root path. (default "/home/federico/Work/dbg-go")
      --target-distro string          target distro to work against. By default tool will work on any supported distro. Can be a regex.
                                      Supported: [almalinux,amazonlinux,amazonlinux2,amazonlinux2022,amazonlinux2023,bottlerocket,centos,debian,fedora,minikube,talos,ubuntu].
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
./dbg-go drivers stats --driver-version 5.0.1+driver
```
</details>

<details>
  <summary>Validate local configs for 5.0.1+driver driver version, for arm64</summary>
  
```bash
./dbg-go configs validate --driver-version 5.0.1+driver --architecture arm64
```
</details>

<details>
  <summary>Generate configs for all supported driver versions by test-infra from kernel-crawler output, for host architecture</summary>
  
```bash
./dbg-go configs generate --repo-root test-infra --auto
```
</details>

<details>
  <summary>Build all x86_64 5.0.1+driver configs, publishing them to s3</summary>
  
```bash
./dbg-go configs build --repo-root test-infra --driver-version 5.0.1+driver --publish
```
</details>

<details>
  <summary>Publish locally built drivers for aarch64 for all supported driver versions by test-infra</summary>

```bash
./dbg-go drivers publish --repo-root test-infra --architecture arm64
```
</details>


> **NOTE:** all commands that require s3 write access, need proper env variables (AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY) exported.

## Bumping driverkit

To bump driverkit, you just need:
```shell
make bump-driverkit DRIVERKIT_VER=vX.Y.Z
```
