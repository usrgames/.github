# Releases

Every usr.games program is released the same way, by one workflow: `.github/workflows/release.yml` in this repository. A repository calls it when a tag `v*` is pushed. This page says what it does, what a repository adds to use it, what a release holds, and how to rehearse.

## What happens on a tag

| Job | What |
|---|---|
| build | Runs the repository's `make check`. Builds an archive for macOS and for Linux, on arm64 and amd64, each holding the program, its `LICENSE`, and its `NOTICES` if there is one. Where asked, builds the program for a web page and renders the manual page to HTML. Writes `manifest.json`. |
| sign | Waits in the repository's `release` environment for a person to approve. Then signs the manifest with the release key, naming its certificate, and checks its own signature. |
| try | On Linux amd64, Linux arm64 and macOS, does what a player's machine does: checks the signature, checks every file against the manifest, unpacks the archive, and runs the program with `--version`. |
| publish | Only once every try has passed. A tag with a hyphen, such as `v1.0.0-rc.1`, is published as a pre-release. |

The keys, and how to check a release by hand, are in [KEYS.md](https://github.com/usrgames/usr.games/blob/main/KEYS.md) of the site's repository.

## What a release holds

| File | What |
|---|---|
| `NAME_OS_ARCH.tar.gz` | The program for one system: `darwin` or `linux`, `arm64` or `amd64`. |
| `NAME.wasm.gz`, `wasm_exec.js` | The program for a web page, and the glue from the Go that built it. Only where asked. |
| `NAME.6.html` | The manual page. Only where asked. |
| `manifest.json` | Every file above, with its size and SHA-256. |
| `manifest.json.sig` | The signature of the manifest. |

### The manifest

```json
{
  "schema": 1,
  "name": "rustwake",
  "version": "v0.1.0",
  "files": [
    {"name": "rustwake.wasm.gz", "size": 987654, "sha256": "..."},
    {"name": "rustwake_darwin_arm64.tar.gz", "size": 1234567, "sha256": "...", "os": "darwin", "arch": "arm64"}
  ]
}
```

- `schema` is 1. A reader refuses a number it does not know.
- `version` is the tag.
- `files` are in name order. An archive also says which `os` and `arch` it runs on; other files do not.
- There is no date in it. The archives are packed with fixed order, dates and owner, so the same tag built twice gives the same manifest.

To trust a release: check the signature of `manifest.json` first, then check every file you use against its size and hash. A file the manifest does not name is not part of the release.

## What a repository adds

One file, `.github/workflows/release.yml`:

```yaml
name: release

on:
  push:
    tags: ["v*"]
  workflow_dispatch:
    inputs:
      rehearsal:
        description: "Try the whole release with keys made for the run. Nothing is published."
        type: boolean
        default: false
      again:
        description: "Run on a tag: sign that release again, under the current certificate."
        type: boolean
        default: false

permissions:
  contents: read

jobs:
  release:
    uses: usrgames/.github/.github/workflows/release.yml@COMMIT
    permissions:
      contents: write
    with:
      name: rustwake
      page: true
      man: man/rustwake.6
      rehearsal: ${{ inputs.rehearsal || false }}
      again: ${{ inputs.again || false }}
    secrets:
      RELEASE_KEY: ${{ secrets.RELEASE_KEY }}
```

`COMMIT` is the full commit of this repository to use. The organization accepts nothing less than a full commit.

And the repository needs:

- **`make check`**, which the build runs first.
- **A `version`.** The program's main package declares `var version = "dev"`, and `NAME --version` prints it. The build sets it to the tag.
- **A `LICENSE`** at the root of the module, and a `NOTICES` if it ships other people's code.
- **A `release` environment** that lets only a tag `v*` in, asks for a reviewer, and holds the release key as the secret `RELEASE_KEY`.

The inputs:

| Input | What | Default |
|---|---|---|
| `name` | The program's name: its file in each archive, and the start of each archive's name. | required |
| `package` | The Go package that is the program. | `.` |
| `dir` | The folder of the Go module, when it is not the repository's root. | `.` |
| `go-version` | The Go to build with. | `1.27` |
| `page` | Also build for a web page. | `false` |
| `man` | A manual page to ship as HTML too. | none |
| `rehearsal` | Sign with keys made for the run, and publish nothing. | `false` |
| `again` | Build nothing: sign this tag's release again. | `false` |

## Rehearse

A rehearsal runs everything but the real key and the publishing. It makes a root, a release key and a certificate inside the run, good for an hour, signs with those, and tries the result on all three systems. It also builds the release twice and fails if the two differ.

- In this repository, every pull request rehearses on the small program in `testdata/hello`.
- In a repository that calls the workflow: run its release workflow by hand, on any branch, with `rehearsal` ticked.

## Sign a release again

A signature carries the release key's certificate, and the certificate ends after a year. On that day a release stops checking out, though nothing in it changed. So at each renewal, a program's latest release is signed again:

1. Renew the certificate, as `KEYS.md` says, and merge it into the site's repository.
2. In each repository, run the release workflow by hand **on the release's tag**, with `again` ticked, and approve it.

The run takes the release as published. Its old signature must have been good on the day it was published, and its manifest must name this program and this tag. It then signs the same manifest under the new certificate, tries the release on all three systems, and only then replaces the signature. A program released since the renewal needs none of this.
