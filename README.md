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

One workflow covers both modes through a `workflow_dispatch` input:

```yaml
name: modernize
on:
  workflow_dispatch:
    inputs:
      mode:
        description: "Execution mode"
        required: true
        default: "fix"
        type: choice
        options:
          - fix
          - check
permissions:
  contents: write
  pull-requests: write
jobs:
  modernize:
    runs-on: ubuntu-latest
    steps:
      - uses: alesr/gomodern@v0.2.0
        with:
          mode: ${{ inputs.mode }}
          create-pr: "true"
```

`check` fails the build if it finds legacy code. `fix` rewrites it and, with
`create-pr: "true"`, opens a PR. The action checks out the repository itself, so
you don't need a separate `actions/checkout` step.

To gate pull requests instead, run on `pull_request` with `mode: check` and
`permissions: { contents: read }`.

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

Enable `Allow GitHub Actions to create and approve pull requests` in the repository action's general settings.

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
