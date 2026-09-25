# gomodern

A GitHub Action that keeps your Go code modern. It runs the official
[`modernize`](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/modernize)
analyzers from `golang.org/x/tools`, then either fails your build when it finds
outdated code or fixes it for you and opens a PR.

It covers the usual suspects: `interface{}` to `any`, hand-written `min`/`max`,
manual slice and map loops that `slices`/`maps` can replace, `sort.Slice` to
`slices.Sort`, and more. Rewrites are gated to the Go version in your `go.mod`,
so it won't push you onto an API your toolchain doesn't have yet.

[![Tests](https://img.shields.io/github/actions/workflow/status/alesr/gomodern/integration-tests.yml?label=tests)](https://github.com/alesr/gomodern/actions)

## Modes

- **`check`** — fail the build if anything needs modernizing, and print the diff.
- **`fix`** — rewrite in place, and optionally open a PR with the changes.

## Usage

### Check mode

Drop this into a PR workflow to block legacy idioms from landing:

```yaml
name: ci
on: [pull_request]
permissions:
  contents: read
jobs:
  modernize:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: alesr/gomodern@v1
        with:
          mode: check
```

### Fix mode

Run it on a schedule to keep the codebase moving:

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
      - uses: actions/checkout@v4
      - uses: alesr/gomodern@v1
        with:
          mode: fix
          create-pr: "true"
          github-token: ${{ secrets.GITHUB_TOKEN }}
```

## Inputs

| Input | Type | Default | Description |
| --- | --- | --- | --- |
| `mode` | `string` | `check` | `check` or `fix` |
| `packages` | `string` | `./...` | Package pattern to analyze |
| `flags` | `string` | `""` | Extra flags for `modernize` (e.g. `-any=true`) |
| `go-version` | `string` | `""` | Override the Go version used for gating (defaults to `go.mod`) |
| `create-pr` | `boolean` | `false` | In `fix` mode, commit and open a PR |
| `pr-title` | `string` | `refactor(go): modernize codebase syntax` | PR title |
| `pr-branch` | `string` | `gomodern/auto-fix` | PR branch |
| `commit-message` | `string` | `refactor: apply gomodern AST rewrites` | Commit message |
| `github-token` | `string` | `""` | Token for PR creation (defaults to the workflow token) |

## Outputs

| Output | Description |
| --- | --- |
| `has-changes` | `true` if anything was found or fixed |
| `modified-files` | Space-separated list of files changed in `fix` mode |
| `pr-number` | PR number when `create-pr: true`, otherwise empty |

## Go version gating

The action reads the `go` directive from your `go.mod` and disables any rewrite
that needs a newer Go. On a `go 1.20` module it won't touch anything that relies
on 1.21+ (like `slices` or `min`/`max`). Set `go-version` to force a specific
version instead.

## Permissions

`check` mode only needs `contents: read`. If you use `fix` with `create-pr`,
you'll also need:

```yaml
permissions:
  contents: write
  pull-requests: write
```

## Pinning

Pin the action (and its transitive actions) to a commit SHA rather than a tag in
production:

```yaml
uses: alesr/gomodern@a1b2c3d4e5f6
```

## Local development

Fixtures live under `testdata/` and the integration matrix is in
`.github/workflows/integration-tests.yml`. Workflow and action linting
(`actionlint`, `zizmor`, `action-validator`) runs in `.github/workflows/lint.yml`.

To run the analyzer yourself:

```bash
go install golang.org/x/tools/go/analysis/passes/modernize/cmd/modernize@latest
cd testdata
modernize ./...        # check: exits non-zero on legacy code
modernize -fix ./...   # fix: rewrites in place
```
