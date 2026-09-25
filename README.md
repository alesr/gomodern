# gomodern

A GitHub Action that runs the
[`modernize`](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/modernize)
analyzers from `golang.org/x/tools` against your repository.

[![Tests](https://img.shields.io/github/actions/workflow/status/alesr/gomodern/integration-tests.yml?label=tests)](https://github.com/alesr/gomodern/actions)

## What it does

Two modes:

- `check` — fails the build if it finds legacy code, and prints the diff.
- `fix` — rewrites the code in place, and can open a PR with the changes.

The analyzers cover things like `interface{}` to `any`, hand-written `min`/`max`,
manual `slices`/`maps` loops, and `sort.Slice` to `slices.Sort`. Fixes are limited
to the Go version in `go.mod`, so it won't use an API newer than your toolchain.

## Usage

### Check mode

```yaml
name: modernize
on: [pull_request]
permissions:
  contents: read
jobs:
  modernize:
    runs-on: ubuntu-latest
    steps:
      - uses: alesr/gomodern@v1
        with:
          mode: check
```

### Fix mode

```yaml
name: modernize
on:
  schedule:
    - cron: "0 3 * * 1"
permissions:
  contents: write
  pull-requests: write
jobs:
  modernize:
    runs-on: ubuntu-latest
    steps:
      - uses: alesr/gomodern@v1
        with:
          mode: fix
          create-pr: "true"
```

The action checks out the repository itself, so you don't need a separate
`actions/checkout` step.

## Inputs

| Input | Type | Default | Description |
| --- | --- | --- | --- |
| `mode` | string | `check` | `check` or `fix` |
| `packages` | string | `./...` | Package pattern to analyze |
| `flags` | string | `""` | Extra flags passed to `modernize` |
| `go-version` | string | `""` | Go version to gate fixes on (defaults to `go.mod`) |
| `create-pr` | boolean | `false` | In `fix` mode, open a PR with the changes |
| `pr-title` | string | `refactor(go): modernize codebase syntax` | PR title |
| `pr-branch` | string | `gomodern/auto-fix` | PR branch |
| `commit-message` | string | `refactor: apply gomodern AST rewrites` | Commit message |
| `github-token` | string | `""` | Token for PR creation (defaults to the workflow token) |

## Outputs

| Output | Description |
| --- | --- |
| `has-changes` | `true` if anything was found or fixed |
| `modified-files` | Space-separated list of files changed in `fix` mode |
| `pr-number` | PR number when `create-pr: true`, otherwise empty |

## Go version gating

The action reads the `go` directive from `go.mod` and disables any fix that needs
a newer Go version. Set `go-version` to override it.

## Permissions

`check` mode needs `contents: read`. `fix` with `create-pr` also needs:

```yaml
permissions:
  contents: write
  pull-requests: write
```

## Pinning

Pin the action to a commit SHA in production:

```yaml
uses: your-org/gomodern@a1b2c3d4e5f6
```

## Local development

Fixtures are under `testdata/`, and the integration matrix is
`.github/workflows/integration-tests.yml`. Workflow and action linting runs in
`.github/workflows/lint.yml`.

To run the analyzer by hand:

```bash
go run golang.org/x/tools/go/analysis/passes/modernize/cmd/modernize@latest ./...
```

## License

[MIT](LICENSE)
