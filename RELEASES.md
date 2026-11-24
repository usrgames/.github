# Releases

Every usr.games program is released the same way, by one workflow in two halves in this repository: `.github/workflows/build.yml` builds the release, and `.github/workflows/publish.yml` signs and publishes it. A repository calls both when a tag `v*` is pushed. This page says what it does, what a repository adds to use it, what a release holds, and how to rehearse.

## What happens on a tag

| Job | What |
|---|---|
| build | Runs the repository's `make check`, then, where asked, the make target that writes its recording. Builds an archive for macOS and for Linux, on arm64 and amd64, each holding the program, its `LICENSE`, its `NOTICES` if there is one, and its manual page if one is named. Where asked, builds the program for a web page and renders the manual page to HTML. Writes `manifest.json`. |
| sign | Waits in the repository's `release` environment for a person to approve. Then signs the manifest with the release key, naming its certificate, and checks its own signature. |
| try | On Linux amd64, Linux arm64 and macOS, does what a player's machine does: checks the signature, checks every file against the manifest, unpacks the archive, and runs the program with `--version`. |
| publish | Only once every try has passed. A tag with a hyphen, such as `v1.0.0-rc.1`, is published as a pre-release. |

build is the first half, `build.yml`. sign, try and publish are the second, `publish.yml`.

The keys, and how to check a release by hand, are in [KEYS.md](https://github.com/usrgames/usr.games/blob/main/KEYS.md) of the site's repository.

## What a release holds

| File | What |
|---|---|
| `NAME_OS_ARCH.tar.gz` | The program for one system: `darwin` or `linux`, `arm64` or `amd64`. With it, its `LICENSE` and, where one is named, its manual page as it was written, which `tty1 man` opens. |
| `NAME.wasm.gz`, `wasm_exec.js` | The program for a web page, and the glue from the Go that built it. Only where asked. |
| `NAME.6.html` | The manual page as HTML, for the site. Only where asked. |
| `NAME.gif`, `NAME.png` | The recording the site shows on the game's page: a few seconds that play over and over, and its first frame as a still, for people who ask for less motion. Only where asked. |
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

One file, `.github/workflows/release.yml`, which calls both halves. tty1's, for example:

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
  build:
    # Signing a release again builds nothing.
    if: ${{ !inputs.again }}
    uses: usrgames/.github/.github/workflows/build.yml@COMMIT
    with:
      name: tty1
      man: man/tty1.1
      rehearsal: ${{ inputs.rehearsal || false }}

  publish:
    needs: build
    if: ${{ !cancelled() && (needs.build.result == 'success' || inputs.again) }}
    uses: usrgames/.github/.github/workflows/publish.yml@COMMIT
    # The half's last job publishes the release. Every other job there only reads.
    permissions:
      contents: write
    with:
      name: tty1
      man: man/tty1.1
      rehearsal: ${{ inputs.rehearsal || false }}
      again: ${{ inputs.again || false }}
    secrets:
      RELEASE_KEY: ${{ secrets.RELEASE_KEY }}
```

`COMMIT` is the full commit of this repository to use, the same in both. The organization accepts nothing less than a full commit. Each half must be called by the repository's own workflow: a call from one shared workflow to another would be looked up in the calling repository.

And the repository needs:

- **`make check`**, which the build runs first. It must leave nothing behind in the checkout: Go records in a program whether it was built from changed sources, and the build stops if it would be.
- **A `version`.** The program's main package declares `var version = "dev"`, and `NAME --version` prints it. The build sets it to the tag.
- **A `LICENSE`** at the root of the module, and a `NOTICES` if it ships other people's code.
- **A `release` environment** that lets only a tag `v*` in, asks for a reviewer, and holds the release key as the secret `RELEASE_KEY`.
- **For a recording, a make target** that writes `NAME.gif` and `NAME.png`, and nothing else, into the folder `MEDIA` names, such as `make media MEDIA=/some/folder`. The GIF plays over and over and is at most 3,000,000 bytes, since the site's front page loads one for each game. The PNG is its first frame. Made twice, they must be the same bytes, which a rehearsal checks, and the target must write nothing into the checkout.

The inputs:

| Input | Half | What | Default |
|---|---|---|---|
| `name` | both | The program's name: its file in each archive, and the start of each archive's name. | required |
| `package` | build | The Go package that is the program. | `.` |
| `dir` | build | The folder of the Go module, when it is not the repository's root. | `.` |
| `go-version` | build | The Go to build with. | `1.27` |
| `page` | build | Also build for a web page. | `false` |
| `man` | both | A manual page: build puts it in each archive and ships it as HTML too, and publish checks each archive holds it. | none |
| `media` | build | A make target that writes the recording, `NAME.gif` and `NAME.png`, into the folder `MEDIA` names. | none |
| `rehearsal` | both | Sign with keys made for the run, and publish nothing. | `false` |
| `again` | publish | Sign this tag's published release again. The caller skips build. | `false` |

## A program whose source is private

Its source stays in its private repository, and its releases are published from a public one beside it, named for its manual page: `usrgames/NAME.6`. A private repository on GitHub Free has no environment, so nothing there could ask a person before the release key signs, and tty1 and the site download a release without signing in.

| Where | What happens |
|---|---|
| The private repository, on a tag `v*` | build, then hand-over: the unsigned release goes to the release repository as a draft, on a tag of the same name, and that repository's release workflow starts. |
| The release repository | take holds the draft to its manifest, then sign waits for a person, then try, then publish, which publishes the bytes that were tried and deletes the draft. |

The private repository's `.github/workflows/release.yml`:

```yaml
name: release

on:
  push:
    tags: ["v*"]
  workflow_dispatch:
    inputs:
      rehearsal:
        description: "Try the whole release with keys made for the run. Nothing is handed over."
        type: boolean
        default: false

permissions:
  contents: read

jobs:
  build:
    uses: usrgames/.github/.github/workflows/build.yml@COMMIT
    with:
      name: rustwake
      page: true
      man: man/rustwake.6
      media: media
      to: usrgames/rustwake.6
      rehearsal: ${{ inputs.rehearsal || false }}
    secrets:
      HANDOVER_TOKEN: ${{ secrets.HANDOVER_TOKEN }}

  # Only a rehearsal signs and tries here, with keys made for the run. It costs this repository's
  # minutes, and a macOS minute counts as ten.
  rehearse:
    if: ${{ inputs.rehearsal }}
    needs: build
    uses: usrgames/.github/.github/workflows/publish.yml@COMMIT
    permissions:
      contents: write
    with:
      name: rustwake
      man: man/rustwake.6
      rehearsal: true
```

The release repository's `.github/workflows/release.yml`. The hand-over starts it by that name:

```yaml
name: release

on:
  workflow_dispatch:
    inputs:
      again:
        description: "Run on a tag: sign that release again, under the current certificate."
        type: boolean
        default: false

permissions:
  contents: read

jobs:
  publish:
    uses: usrgames/.github/.github/workflows/publish.yml@COMMIT
    permissions:
      contents: write
    with:
      name: rustwake
      man: man/rustwake.6
      draft: ${{ !inputs.again }}
      again: ${{ inputs.again }}
    secrets:
      RELEASE_KEY: ${{ secrets.RELEASE_KEY }}
```

The private repository needs what any repository does, above, except the `release` environment and the key. It also needs:

- **The secret `HANDOVER_TOKEN`:** a fine-grained token whose owner is `usrgames`, for the release repositories only, that may read and write their contents and their actions. It ends with the release key's certificate and is made again at each renewal, as `KEYS.md` says.

The release repository needs:

- **To be public, and to hold none of the source:** its `README.md`, a `LICENSE` that says what a player may do with the program, its `NOTICES`, and the workflow above.
- **A rule on `main` and on tags `v*`.** The owner of the hand-over token makes the tags, so may bypass it.
- **A `release` environment** that lets only a tag `v*` in, asks for a reviewer, and holds the release key as the secret `RELEASE_KEY`.
- **Issues on**, for players.

To hand a version over again, after a take or a try that failed or a release that was not approved: delete its draft release and its tag in the release repository, which is safe since nothing is signed yet, then run the release again in the private repository. A run in the release repository that failed after take can be re-run there as it is.

### What a closed program still shows

Its source stays private, but a program is not a secret. Anyone with a release can read in it the names of its functions, types and packages, the modules it is built from and their versions, the commit and the time it was built from (`go version -m`), and every string in it. The build for a web page is a download like the others. The builds leave out the paths of the machine that made them and the tables a debugger reads (`-trimpath`, `-s -w`).

## Rehearse

A rehearsal runs everything but the real key and the publishing. It makes a root, a release key and a certificate inside the run, good for an hour, signs with those, and tries the result on all three systems. It also builds the release twice and fails if the two differ.

- In this repository, every pull request rehearses both halves on the small program in `testdata/hello`.
- In a repository that calls the workflow: run its release workflow by hand, on any branch, with `rehearsal` ticked.

## Sign a release again

A signature carries the release key's certificate, and the certificate ends after a year. On that day a release stops checking out, though nothing in it changed. So at each renewal, a program's latest release is signed again:

1. Renew the certificate, as `KEYS.md` says, and merge it into the site's repository.
2. In each repository, run the release workflow by hand **on the release's tag**, with `again` ticked, and approve it. For a program whose source is private, that is its release repository's workflow.

The run takes the release as published. Its old signature must have been good on the day it was published, and its manifest must name this program and this tag. It then signs the same manifest under the new certificate, tries the release on all three systems, and only then replaces the signature. A program released since the renewal needs none of this.
