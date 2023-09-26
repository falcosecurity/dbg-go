# Release Process

Our release process is fully automated using [Github actions](.github/workflows/release.yml) and [goreleaser](https://github.com/goreleaser/goreleaser) tool for artifacts.

When we release we do the following process:

1. We decide together (usually in the #falco channel in [slack](https://kubernetes.slack.com/messages/falco)) what's the next version to tag
2. A person with repository rights does the tag
3. The same person runs commands in their machine following the "Release commands" section below
4. Once the CI has done its job, the tag is live on [Github](https://github.com/falcosecurity/dbg-go/releases) with the artifacts, and the container image is live on [GHCR](ghcr.io/falcosecurity/dbg-go).

## Release commands

Tag the version

```bash
git tag -a v0.1.0-rc.0 -m "v0.1.0-rc.0"
git push origin v0.1.0-rc.0
```
