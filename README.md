# .github

What every usr.games repository shares.

| Path | What |
|---|---|
| `profile/README.md` | The page GitHub shows at the top of the organization. |
| `.github/workflows/release.yml` | The one workflow every repository releases with. [RELEASES.md](RELEASES.md) says what it does and how to call it. |
| `.github/workflows/ci.yml` | Rehearses that workflow on every pull request here. |
| `testdata/hello` | The small program the rehearsal releases. |

The picture at the top of the profile page is `site/card.png` of [usr.games](https://github.com/usrgames/usr.games), loaded from its `main` branch. A test draws it there, so it is not copied here.

When a repository joins the organization, or changes what it holds, its row on the profile page changes too.

## License

MIT. See [LICENSE](LICENSE).
