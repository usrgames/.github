# .github

What every usr.games repository shares.

| Path | What |
|---|---|
| `profile/README.md` | The page GitHub shows at the top of the organization. |
| `.github/workflows/build.yml`, `.github/workflows/publish.yml` | The two halves of the one workflow every repository releases with: the first builds a release, the second signs and publishes it. [RELEASES.md](RELEASES.md) says what they do and how to call them. |
| `.github/workflows/ci.yml` | Rehearses both halves on every pull request here. |
| `testdata/hello` | The small program the rehearsal releases. |

The picture at the top of the profile page is `site/card.png` of [usr.games](https://github.com/usrgames/usr.games), loaded from its `main` branch. A test draws it there, so it is not copied here.

When a repository joins the organization, or changes what it holds, its row on the profile page changes too.

## License

MIT. See [LICENSE](LICENSE).
