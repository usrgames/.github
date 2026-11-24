# .github

What every usr.games repository shares.

| Path | What |
|---|---|
| `profile/README.md` | The page GitHub shows at the top of the organization. |
| `.github/workflows/build.yml`, `.github/workflows/publish.yml` | The two halves of the one workflow every repository releases with: the first builds a release, the second signs and publishes it. [RELEASES.md](RELEASES.md) says what they do and how to call them. |
| `.github/workflows/retry.yml` | Reruns a run's failed jobs when GitHub found no runner for one of them. [RELEASES.md](RELEASES.md) says how to call it. |
| `.github/workflows/rerun.yml` | Calls it for this repository's `ci`. |
| `.github/workflows/ci.yml` | Rehearses both halves, and the rule in `retry.yml`, on every pull request here. |
| `testdata/hello` | The small program the rehearsal releases. |
| `testdata/retry` | The runs the rehearsal of `retry.yml` reads, saved from GitHub's API. |

The picture at the top of the profile page is `site/card.png` of [usr.games](https://github.com/usrgames/usr.games), loaded from its `main` branch. A test draws it there, so it is not copied here.

The profile page's table lists only what a visitor can open: a private repository has no row, and a game's release repository gets one when it is made. When a public repository joins the organization, or changes what it holds, its row changes too.

## License

MIT. See [LICENSE](LICENSE).
